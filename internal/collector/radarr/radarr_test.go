package radarr

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/st0o0/tentacle/internal/client/arr"
)

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v3/system/status":
			w.Write([]byte(`{"version":"5.3.6","branch":"main","runtimeName":"docker","runtimeVersion":"6.0.0","startTime":"2024-01-01T00:00:00Z"}`))
		case "/api/v3/health":
			w.Write([]byte(`[{"source":"IndexerStatusCheck","type":"warning","message":"no indexers","wikiUrl":"http://wiki"}]`))
		case "/api/v3/movie":
			w.Write([]byte(`[
				{"title":"Movie A","monitored":true,"hasFile":true,"sizeOnDisk":1500000000,"status":"released","year":2023},
				{"title":"Movie B","monitored":true,"hasFile":false,"sizeOnDisk":0,"status":"announced","year":2024},
				{"title":"Movie C","monitored":false,"hasFile":true,"sizeOnDisk":2500000000,"status":"released","year":2022}
			]`))
		case "/api/v3/wanted/missing":
			w.Write([]byte(`{"totalRecords":5}`))
		case "/api/v3/queue":
			w.Write([]byte(`{"totalRecords":3,"records":[]}`))
		case "/api/v3/rootfolder":
			w.Write([]byte(`[
				{"path":"/movies","freeSpace":500000000000,"totalSpace":1000000000000},
				{"path":"/movies2","freeSpace":200000000000,"totalSpace":400000000000}
			]`))
		case "/api/v3/calendar":
			w.Write([]byte(`[{"title":"Upcoming A"},{"title":"Upcoming B"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
}

func newClient(url string) *arr.Client {
	return arr.NewClient(url, "test-api-key", http.DefaultClient)
}

// gather collects metrics from a collector and returns them keyed by fqName.
// For metrics with labels, the key is fqName::labelName=labelValue.
func gather(t *testing.T, c prometheus.Collector) map[string]float64 {
	t.Helper()
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	result := map[string]float64{}
	for _, mf := range families {
		for _, m := range mf.GetMetric() {
			key := mf.GetName()
			for _, l := range m.GetLabel() {
				key += "::" + l.GetName() + "=" + l.GetValue()
			}
			result[key] = value(m)
		}
	}
	return result
}

func value(m *dto.Metric) float64 {
	if g := m.GetGauge(); g != nil {
		return g.GetValue()
	}
	return 0
}

func assertMetric(t *testing.T, metrics map[string]float64, key string, want float64) {
	t.Helper()
	got, ok := metrics[key]
	if !ok {
		t.Fatalf("metric %q not found", key)
	}
	if got != want {
		t.Fatalf("metric %q = %v, want %v", key, got, want)
	}
}

func TestSystemCollector(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	m := gather(t, NewSystemCollector(newClient(srv.URL), 5*time.Second, slog.Default()))

	assertMetric(t, m, "radarr_up", 1)
	assertMetric(t, m, "radarr_system_info::branch=main::runtime=docker::version=5.3.6", 1)
	assertMetric(t, m, "radarr_health_issues_total::source=IndexerStatusCheck::type=warning", 1)
	assertMetric(t, m, "radarr_scrape_success::collector=system", 1)
}

func TestMoviesCollector(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	m := gather(t, NewMoviesCollector(newClient(srv.URL), 5*time.Second, slog.Default()))

	assertMetric(t, m, "radarr_movies_total", 3)
	assertMetric(t, m, "radarr_movies_monitored_total", 2)
	assertMetric(t, m, "radarr_movies_downloaded_total", 2)
	assertMetric(t, m, "radarr_movies_size_bytes", 4000000000)
	assertMetric(t, m, "radarr_movies_missing_total", 5)
	assertMetric(t, m, "radarr_scrape_success::collector=movies", 1)
}

func TestQueueCollector(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	m := gather(t, NewQueueCollector(newClient(srv.URL), 5*time.Second, slog.Default()))

	assertMetric(t, m, "radarr_queue_total", 3)
	assertMetric(t, m, "radarr_scrape_success::collector=queue", 1)
}

func TestDiskCollector(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	m := gather(t, NewDiskCollector(newClient(srv.URL), 5*time.Second, slog.Default()))

	assertMetric(t, m, "radarr_disk_total_bytes::path=/movies", 1000000000000)
	assertMetric(t, m, "radarr_disk_free_bytes::path=/movies", 500000000000)
	assertMetric(t, m, "radarr_disk_total_bytes::path=/movies2", 400000000000)
	assertMetric(t, m, "radarr_disk_free_bytes::path=/movies2", 200000000000)
	assertMetric(t, m, "radarr_scrape_success::collector=disk", 1)
}

func TestCalendarCollector(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	m := gather(t, NewCalendarCollector(newClient(srv.URL), 5*time.Second, slog.Default()))

	assertMetric(t, m, "radarr_calendar_upcoming_total", 2)
	assertMetric(t, m, "radarr_scrape_success::collector=calendar", 1)
}
