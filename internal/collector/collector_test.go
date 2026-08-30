package collector

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type stubSub struct {
	name    string
	desc    *prometheus.Desc
	value   float64
	failing bool
}

func (s *stubSub) Name() string { return s.name }

func (s *stubSub) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.desc
}

func (s *stubSub) Update(ch chan<- prometheus.Metric) error {
	if s.failing {
		return fmt.Errorf("stub error")
	}
	ch <- prometheus.MustNewConstMetric(s.desc, prometheus.GaugeValue, s.value)
	return nil
}

func newStub(name string, value float64, failing bool) *stubSub {
	return &stubSub{
		name:    name,
		desc:    prometheus.NewDesc("test_"+name+"_value", "stub", nil, nil),
		value:   value,
		failing: failing,
	}
}

func gather(t *testing.T, c prometheus.Collector) map[string]*dto.MetricFamily {
	t.Helper()
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(c)
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	result := make(map[string]*dto.MetricFamily, len(families))
	for _, mf := range families {
		result[mf.GetName()] = mf
	}
	return result
}

func findGauge(families map[string]*dto.MetricFamily, name string, labels map[string]string) *float64 {
	mf, ok := families[name]
	if !ok {
		return nil
	}
	for _, m := range mf.GetMetric() {
		actual := make(map[string]string, len(m.GetLabel()))
		for _, lp := range m.GetLabel() {
			actual[lp.GetName()] = lp.GetValue()
		}
		match := true
		for k, v := range labels {
			if actual[k] != v {
				match = false
				break
			}
		}
		if match {
			v := m.GetGauge().GetValue()
			return &v
		}
	}
	return nil
}

func assertGauge(t *testing.T, families map[string]*dto.MetricFamily, name string, labels map[string]string, expected float64) {
	t.Helper()
	v := findGauge(families, name, labels)
	if v == nil {
		t.Errorf("metric %s with labels %v not found", name, labels)
		return
	}
	if *v != expected {
		t.Errorf("metric %s%v = %v, want %v", name, labels, *v, expected)
	}
}

func TestServiceCollector_NoDuplicateDescriptorPanic(t *testing.T) {
	sc := NewServiceCollector("test", slog.Default(),
		newStub("alpha", 1, false),
		newStub("beta", 2, false),
	)
	families := gather(t, sc)

	assertGauge(t, families, "test_alpha_value", nil, 1)
	assertGauge(t, families, "test_beta_value", nil, 2)
	assertGauge(t, families, "test_scrape_success", map[string]string{"collector": "alpha"}, 1)
	assertGauge(t, families, "test_scrape_success", map[string]string{"collector": "beta"}, 1)

	if findGauge(families, "test_scrape_duration_seconds", map[string]string{"collector": "alpha"}) == nil {
		t.Error("missing scrape_duration_seconds for alpha")
	}
	if findGauge(families, "test_scrape_duration_seconds", map[string]string{"collector": "beta"}) == nil {
		t.Error("missing scrape_duration_seconds for beta")
	}
}

func TestServiceCollector_FailureIsolation(t *testing.T) {
	sc := NewServiceCollector("test", slog.Default(),
		newStub("good", 42, false),
		newStub("bad", 0, true),
	)
	families := gather(t, sc)

	assertGauge(t, families, "test_good_value", nil, 42)
	assertGauge(t, families, "test_scrape_success", map[string]string{"collector": "good"}, 1)
	assertGauge(t, families, "test_scrape_success", map[string]string{"collector": "bad"}, 0)
}

func TestServiceCollector_Empty(t *testing.T) {
	sc := NewServiceCollector("test", slog.Default())
	families := gather(t, sc)

	if len(families) != 0 {
		t.Errorf("expected no metrics, got %d families", len(families))
	}
}

func TestCachedCollector_ServesCachedMetrics(t *testing.T) {
	stub := newStub("cached", 99, false)
	cc := NewCachedCollector(stub, time.Hour, slog.Default())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cc.Start(ctx)

	// Wait for the initial refresh to complete.
	time.Sleep(50 * time.Millisecond)

	sc := NewServiceCollector("test", slog.Default(), cc)
	families := gather(t, sc)
	assertGauge(t, families, "test_cached_value", nil, 99)
	assertGauge(t, families, "test_scrape_success", map[string]string{"collector": "cached"}, 1)
}

func TestCachedCollector_EmptyBeforeStart(t *testing.T) {
	stub := newStub("lazy", 1, false)
	cc := NewCachedCollector(stub, time.Hour, slog.Default())

	ch := make(chan prometheus.Metric, 16)
	err := cc.Update(ch)
	close(ch)

	if err != nil {
		t.Errorf("expected nil error before start, got %v", err)
	}
	var count int
	for range ch {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 metrics before start, got %d", count)
	}
}

func TestCachedCollector_PropagatesError(t *testing.T) {
	stub := newStub("failing", 0, true)
	cc := NewCachedCollector(stub, time.Hour, slog.Default())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cc.Start(ctx)

	time.Sleep(50 * time.Millisecond)

	ch := make(chan prometheus.Metric, 16)
	err := cc.Update(ch)
	close(ch)

	if err == nil {
		t.Error("expected error from cached failing collector")
	}
}

func TestCachedCollector_UpdatesOnTick(t *testing.T) {
	stub := newStub("ticking", 10, false)
	cc := NewCachedCollector(stub, 50*time.Millisecond, slog.Default())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cc.Start(ctx)

	time.Sleep(30 * time.Millisecond)

	stub.value = 20
	time.Sleep(80 * time.Millisecond)

	sc := NewServiceCollector("test", slog.Default(), cc)
	families := gather(t, sc)
	assertGauge(t, families, "test_ticking_value", nil, 20)
}
