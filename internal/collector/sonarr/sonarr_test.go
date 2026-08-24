package sonarr

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/st0o0/tentacle/internal/client/arr"
)

// newTestServer creates an httptest.Server that routes API requests to canned JSON responses.
func newTestServer(t *testing.T, responses map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip query parameters for matching; match on path only.
		path := r.URL.Path
		data, ok := responses[path]
		if !ok {
			t.Errorf("unexpected request path: %s", path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			t.Fatalf("failed to encode response for %s: %v", path, err)
		}
	}))
}

func newClient(serverURL string) *arr.Client {
	return arr.NewClient(serverURL, "test-api-key", http.DefaultClient)
}

func TestSystemCollector(t *testing.T) {
	srv := newTestServer(t, map[string]any{
		"/api/v3/system/status": arr.SystemStatus{
			Version:     "4.0.0.1",
			Branch:      "main",
			RuntimeName: "dotnet",
			StartTime:   "2024-01-01T00:00:00Z",
		},
		"/api/v3/health": []arr.HealthCheck{
			{Type: "warning", Source: "IndexerStatusCheck"},
			{Type: "error", Source: "DownloadClientCheck"},
		},
	})
	defer srv.Close()

	c := NewSystemCollector(newClient(srv.URL), 5*time.Second, slog.Default())

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}

	found := map[string]bool{
		"sonarr_up":                           false,
		"sonarr_system_info":                  false,
		"sonarr_system_start_time_seconds":    false,
		"sonarr_health_issues_total":          false,
	}

	for _, mf := range mfs {
		name := mf.GetName()
		switch name {
		case "sonarr_up":
			found[name] = true
			if v := mf.GetMetric()[0].GetGauge().GetValue(); v != 1 {
				t.Errorf("expected sonarr_up=1, got %v", v)
			}
		case "sonarr_system_info":
			found[name] = true
			m := mf.GetMetric()[0]
			labels := map[string]string{}
			for _, lp := range m.GetLabel() {
				labels[lp.GetName()] = lp.GetValue()
			}
			if labels["version"] != "4.0.0.1" {
				t.Errorf("expected version=4.0.0.1, got %s", labels["version"])
			}
			if labels["branch"] != "main" {
				t.Errorf("expected branch=main, got %s", labels["branch"])
			}
			if labels["runtime"] != "dotnet" {
				t.Errorf("expected runtime=dotnet, got %s", labels["runtime"])
			}
		case "sonarr_system_start_time_seconds":
			found[name] = true
			if v := mf.GetMetric()[0].GetGauge().GetValue(); v != 1704067200 {
				t.Errorf("expected start_time=1704067200, got %v", v)
			}
		case "sonarr_health_issues_total":
			found[name] = true
			if len(mf.GetMetric()) != 2 {
				t.Errorf("expected 2 health issue metrics, got %d", len(mf.GetMetric()))
			}
		}
	}

	for name, ok := range found {
		if !ok {
			t.Errorf("metric %s not found in gathered output", name)
		}
	}
}

func TestSeriesCollector(t *testing.T) {
	srv := newTestServer(t, map[string]any{
		"/api/v3/series": []arr.Series{
			{Title: "Show A", Monitored: true, Status: "continuing", SeasonCount: 5, TotalEpisodeCount: 20, EpisodeFileCount: 15, SizeOnDisk: 5000000000},
			{Title: "Show B", Monitored: false, Status: "ended", SeasonCount: 3, TotalEpisodeCount: 10, EpisodeFileCount: 10, SizeOnDisk: 3000000000},
			{Title: "Show C", Monitored: true, Status: "continuing", SeasonCount: 2, TotalEpisodeCount: 5, EpisodeFileCount: 2, SizeOnDisk: 1000000000},
		},
		"/api/v3/wanted/missing": arr.WantedResponse{TotalRecords: 8},
	})
	defer srv.Close()

	c := NewSeriesCollector(newClient(srv.URL), 5*time.Second, slog.Default())

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}

	expected := map[string]float64{
		"sonarr_series_total":              3,
		"sonarr_series_monitored_total":    2,
		"sonarr_seasons_total":             10, // 5+3+2
		"sonarr_episodes_total":            35, // 20+10+5
		"sonarr_episodes_downloaded_total": 27, // 15+10+2
		"sonarr_episodes_missing_total":    8,
		"sonarr_series_size_bytes":         9000000000, // 5B+3B+1B
	}

	for _, mf := range mfs {
		name := mf.GetName()
		if exp, ok := expected[name]; ok {
			got := mf.GetMetric()[0].GetGauge().GetValue()
			if got != exp {
				t.Errorf("%s: expected %v, got %v", name, exp, got)
			}
			delete(expected, name)
		}
	}

	for name := range expected {
		t.Errorf("metric %s not found", name)
	}
}

func TestQueueCollector(t *testing.T) {
	srv := newTestServer(t, map[string]any{
		"/api/v3/queue": arr.QueueResponse{
			TotalRecords: 4,
			Records: []arr.QueueRecord{
				{TrackedDownloadState: "downloading"},
				{TrackedDownloadState: "downloading"},
				{TrackedDownloadState: "importPending"},
				{TrackedDownloadState: "failedPending"},
			},
		},
	})
	defer srv.Close()

	c := NewQueueCollector(newClient(srv.URL), 5*time.Second, slog.Default())

	expected := strings.NewReader(`
# HELP sonarr_queue_total Total number of items in the download queue.
# TYPE sonarr_queue_total gauge
sonarr_queue_total 4
`)
	if err := testutil.CollectAndCompare(c, expected, "sonarr_queue_total"); err != nil {
		t.Errorf("queue total mismatch: %v", err)
	}

	// Verify by_state counts exist
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}
	found := false
	for _, mf := range mfs {
		if mf.GetName() == "sonarr_queue_by_state_total" {
			found = true
			if len(mf.GetMetric()) != 3 {
				t.Errorf("expected 3 state metrics, got %d", len(mf.GetMetric()))
			}
		}
	}
	if !found {
		t.Error("sonarr_queue_by_state_total not found")
	}
}

func TestExtrasCollector(t *testing.T) {
	srv := newTestServer(t, map[string]any{
		"/api/v3/system/backup": []arr.Backup{
			{Id: 1, Name: "backup1", Time: "2024-06-15T08:00:00Z"},
			{Id: 2, Name: "backup2", Time: "2024-06-14T08:00:00Z"},
		},
		"/api/v3/update": []arr.Update{
			{Version: "4.1.0", Installed: false, Latest: true},
			{Version: "4.0.0", Installed: true, Latest: false},
		},
		"/api/v3/blocklist": arr.BlocklistResponse{TotalRecords: 15},
		"/api/v3/downloadclient": []arr.DownloadClient{
			{Name: "SABnzbd", Protocol: "usenet", Priority: 1, Enable: true},
			{Name: "Disabled", Protocol: "torrent", Priority: 2, Enable: false},
		},
	})
	defer srv.Close()

	c := NewExtrasCollector(newClient(srv.URL), 5*time.Second, slog.Default())

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}

	checks := map[string]float64{
		"sonarr_backup_total":    2,
		"sonarr_blocklist_total": 15,
	}

	for _, mf := range mfs {
		name := mf.GetName()
		if exp, ok := checks[name]; ok {
			got := mf.GetMetric()[0].GetGauge().GetValue()
			if got != exp {
				t.Errorf("%s: expected %v, got %v", name, exp, got)
			}
			delete(checks, name)
		}
		if name == "sonarr_update_available" {
			m := mf.GetMetric()[0]
			if m.GetGauge().GetValue() != 1 {
				t.Errorf("expected update_available=1, got %v", m.GetGauge().GetValue())
			}
		}
		if name == "sonarr_download_client_info" {
			if len(mf.GetMetric()) != 1 {
				t.Errorf("expected 1 download client info (only enabled), got %d", len(mf.GetMetric()))
			}
		}
	}

	for name := range checks {
		t.Errorf("metric %s not found", name)
	}
}

func TestDiskCollector(t *testing.T) {
	srv := newTestServer(t, map[string]any{
		"/api/v3/rootfolder": []arr.RootFolder{
			{Path: "/tv", TotalSpace: 1000000000000, FreeSpace: 250000000000},
			{Path: "/anime", TotalSpace: 500000000000, FreeSpace: 100000000000},
		},
	})
	defer srv.Close()

	c := NewDiskCollector(newClient(srv.URL), 5*time.Second, slog.Default())

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}

	type diskExpect struct {
		path  string
		total float64
		free  float64
	}
	expects := []diskExpect{
		{"/tv", 1000000000000, 250000000000},
		{"/anime", 500000000000, 100000000000},
	}

	totalByPath := map[string]float64{}
	freeByPath := map[string]float64{}

	for _, mf := range mfs {
		for _, m := range mf.GetMetric() {
			var pathVal string
			for _, lp := range m.GetLabel() {
				if lp.GetName() == "path" {
					pathVal = lp.GetValue()
				}
			}
			if pathVal == "" {
				continue
			}
			switch mf.GetName() {
			case "sonarr_disk_total_bytes":
				totalByPath[pathVal] = m.GetGauge().GetValue()
			case "sonarr_disk_free_bytes":
				freeByPath[pathVal] = m.GetGauge().GetValue()
			}
		}
	}

	for _, e := range expects {
		if v, ok := totalByPath[e.path]; !ok || v != e.total {
			t.Errorf("disk total_bytes for %s: expected %v, got %v (found=%v)", e.path, e.total, v, ok)
		}
		if v, ok := freeByPath[e.path]; !ok || v != e.free {
			t.Errorf("disk free_bytes for %s: expected %v, got %v (found=%v)", e.path, e.free, v, ok)
		}
	}
}

func TestCalendarCollector(t *testing.T) {
	srv := newTestServer(t, map[string]any{
		"/api/v3/calendar": []arr.CalendarEntry{
			{Title: "Episode 1"},
			{Title: "Episode 2"},
			{Title: "Episode 3"},
		},
	})
	defer srv.Close()

	c := NewCalendarCollector(newClient(srv.URL), 5*time.Second, slog.Default())

	expected := strings.NewReader(`
# HELP sonarr_calendar_upcoming_total Total number of upcoming episodes in the next 7 days.
# TYPE sonarr_calendar_upcoming_total gauge
sonarr_calendar_upcoming_total 3
`)
	if err := testutil.CollectAndCompare(c, expected, "sonarr_calendar_upcoming_total"); err != nil {
		t.Errorf("calendar metric mismatch: %v", err)
	}
}
