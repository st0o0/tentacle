package prowlarr

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/st0o0/tentacle/internal/client/arr"
)

func newTestServer(mux *http.ServeMux) (*httptest.Server, *arr.Client) {
	srv := httptest.NewServer(mux)
	client := arr.NewClient(srv.URL, "test-api-key", srv.Client())
	return srv, client
}

func TestSystemCollector_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/system/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"version": "1.12.0",
			"branch": "main",
			"runtimeName": ".NET",
			"runtimeVersion": "8.0",
			"startTime": "2024-01-01T00:00:00Z"
		}`))
	})
	mux.HandleFunc("/api/v3/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"source": "IndexerRssCheck", "type": "warning", "message": "msg1", "wikiUrl": ""},
			{"source": "IndexerRssCheck", "type": "warning", "message": "msg2", "wikiUrl": ""},
			{"source": "UpdateCheck", "type": "error", "message": "update", "wikiUrl": ""}
		]`))
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewSystemCollector(client, 5*time.Second, slog.Default())

	// up = 1
	expected := `
		# HELP prowlarr_up Whether Prowlarr is reachable.
		# TYPE prowlarr_up gauge
		prowlarr_up 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_up"); err != nil {
		t.Error(err)
	}

	// system_info with labels
	expected = `
		# HELP prowlarr_system_info Prowlarr system information.
		# TYPE prowlarr_system_info gauge
		prowlarr_system_info{branch="main",runtime=".NET 8.0",version="1.12.0"} 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_system_info"); err != nil {
		t.Error(err)
	}

	// health issues grouped by type+source
	expected = `
		# HELP prowlarr_health_issues_total Number of health issues by type and source.
		# TYPE prowlarr_health_issues_total gauge
		prowlarr_health_issues_total{source="IndexerRssCheck",type="warning"} 2
		prowlarr_health_issues_total{source="UpdateCheck",type="error"} 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_health_issues_total"); err != nil {
		t.Error(err)
	}

	// start_time
	expected = `
		# HELP prowlarr_system_start_time_seconds Unix timestamp of when Prowlarr started.
		# TYPE prowlarr_system_start_time_seconds gauge
		prowlarr_system_start_time_seconds 1.7040672e+09
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_system_start_time_seconds"); err != nil {
		t.Error(err)
	}

	// total metric count: up + info + start_time + 2 health + scrape_duration + scrape_success = 7
	if count := testutil.CollectAndCount(c); count != 7 {
		t.Errorf("metric count = %d, want 7", count)
	}
}

func TestSystemCollector_APIFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/system/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewSystemCollector(client, 5*time.Second, slog.Default())

	expected := `
		# HELP prowlarr_up Whether Prowlarr is reachable.
		# TYPE prowlarr_up gauge
		prowlarr_up 0
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_up"); err != nil {
		t.Error(err)
	}

	// scrape_success = 0 for system
	expected = `
		# HELP prowlarr_scrape_success Whether a collector scrape was successful.
		# TYPE prowlarr_scrape_success gauge
		prowlarr_scrape_success{collector="system"} 0
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_scrape_success"); err != nil {
		t.Error(err)
	}
}

func TestIndexersCollector_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/indexer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"name": "NZBgeek", "enable": true, "protocol": "usenet", "priority": 1},
			{"name": "Torznab", "enable": true, "protocol": "torrent", "priority": 2},
			{"name": "Disabled", "enable": false, "protocol": "torrent", "priority": 3}
		]`))
	})
	mux.HandleFunc("/api/v3/indexerstats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"indexers": [
				{
					"indexerName": "NZBgeek",
					"numberOfQueries": 150,
					"numberOfGrabs": 42,
					"numberOfFailedQueries": 3,
					"numberOfFailedGrabs": 1,
					"averageResponseTime": 500
				},
				{
					"indexerName": "Torznab",
					"numberOfQueries": 200,
					"numberOfGrabs": 10,
					"numberOfFailedQueries": 5,
					"numberOfFailedGrabs": 2,
					"averageResponseTime": 1200
				}
			]
		}`))
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewIndexersCollector(client, 5*time.Second, slog.Default())

	// indexers_total = 3, enabled = 2
	expected := `
		# HELP prowlarr_indexers_total Total number of configured indexers.
		# TYPE prowlarr_indexers_total gauge
		prowlarr_indexers_total 3
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexers_total"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_indexers_enabled_total Total number of enabled indexers.
		# TYPE prowlarr_indexers_enabled_total gauge
		prowlarr_indexers_enabled_total 2
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexers_enabled_total"); err != nil {
		t.Error(err)
	}

	// per-indexer stats
	expected = `
		# HELP prowlarr_indexer_queries_total Total number of queries per indexer.
		# TYPE prowlarr_indexer_queries_total gauge
		prowlarr_indexer_queries_total{indexer="NZBgeek"} 150
		prowlarr_indexer_queries_total{indexer="Torznab"} 200
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexer_queries_total"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_indexer_grabs_total Total number of grabs per indexer.
		# TYPE prowlarr_indexer_grabs_total gauge
		prowlarr_indexer_grabs_total{indexer="NZBgeek"} 42
		prowlarr_indexer_grabs_total{indexer="Torznab"} 10
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexer_grabs_total"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_indexer_failed_queries_total Total number of failed queries per indexer.
		# TYPE prowlarr_indexer_failed_queries_total gauge
		prowlarr_indexer_failed_queries_total{indexer="NZBgeek"} 3
		prowlarr_indexer_failed_queries_total{indexer="Torznab"} 5
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexer_failed_queries_total"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_indexer_failed_grabs_total Total number of failed grabs per indexer.
		# TYPE prowlarr_indexer_failed_grabs_total gauge
		prowlarr_indexer_failed_grabs_total{indexer="NZBgeek"} 1
		prowlarr_indexer_failed_grabs_total{indexer="Torznab"} 2
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexer_failed_grabs_total"); err != nil {
		t.Error(err)
	}

	// avg_response_seconds: 500ms -> 0.5s, 1200ms -> 1.2s
	expected = `
		# HELP prowlarr_indexer_avg_response_seconds Average response time per indexer in seconds.
		# TYPE prowlarr_indexer_avg_response_seconds gauge
		prowlarr_indexer_avg_response_seconds{indexer="NZBgeek"} 0.5
		prowlarr_indexer_avg_response_seconds{indexer="Torznab"} 1.2
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexer_avg_response_seconds"); err != nil {
		t.Error(err)
	}

	// by_protocol
	expected = `
		# HELP prowlarr_indexers_by_protocol_total Number of indexers by protocol.
		# TYPE prowlarr_indexers_by_protocol_total gauge
		prowlarr_indexers_by_protocol_total{protocol="torrent"} 2
		prowlarr_indexers_by_protocol_total{protocol="usenet"} 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexers_by_protocol_total"); err != nil {
		t.Error(err)
	}

	// total: indexersTotal(1) + enabledTotal(1) + byProtocol(2) + 2*5 stats(10) + scrape_duration(1) + scrape_success(1) = 16
	if count := testutil.CollectAndCount(c); count != 16 {
		t.Errorf("metric count = %d, want 16", count)
	}
}

func TestIndexersCollector_APIFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/indexer", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewIndexersCollector(client, 5*time.Second, slog.Default())

	expected := `
		# HELP prowlarr_scrape_success Whether a collector scrape was successful.
		# TYPE prowlarr_scrape_success gauge
		prowlarr_scrape_success{collector="indexers"} 0
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_scrape_success"); err != nil {
		t.Error(err)
	}

	// Only scrape_duration + scrape_success emitted on failure
	if count := testutil.CollectAndCount(c); count != 2 {
		t.Errorf("metric count on failure = %d, want 2", count)
	}
}

func TestIndexersCollector_StatsFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/indexer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"name": "Test", "enable": true, "protocol": "usenet", "priority": 1}]`))
	})
	mux.HandleFunc("/api/v3/indexerstats", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewIndexersCollector(client, 5*time.Second, slog.Default())

	// indexers_total and enabled_total should still be emitted
	expected := `
		# HELP prowlarr_indexers_total Total number of configured indexers.
		# TYPE prowlarr_indexers_total gauge
		prowlarr_indexers_total 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexers_total"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_scrape_success Whether a collector scrape was successful.
		# TYPE prowlarr_scrape_success gauge
		prowlarr_scrape_success{collector="indexers"} 0
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_scrape_success"); err != nil {
		t.Error(err)
	}
}

func TestAppsCollector_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/applications", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"name": "Sonarr", "syncLevel": "fullSync", "implementation": "Sonarr"},
			{"name": "Radarr", "syncLevel": "fullSync", "implementation": "Radarr"}
		]`))
	})
	mux.HandleFunc("/api/v3/indexer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"id": 1, "name": "NZBgeek", "enable": true, "protocol": "usenet", "priority": 1},
			{"id": 2, "name": "Torznab", "enable": true, "protocol": "torrent", "priority": 2}
		]`))
	})
	mux.HandleFunc("/api/v3/indexerstatus", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"indexerId": 1, "disabledTill": "2099-01-01T00:00:00Z"}]`))
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewAppsCollector(client, 5*time.Second, slog.Default())

	expected := `
		# HELP prowlarr_apps_total Total number of connected applications.
		# TYPE prowlarr_apps_total gauge
		prowlarr_apps_total 2
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_apps_total"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_app_info Connected application information.
		# TYPE prowlarr_app_info gauge
		prowlarr_app_info{implementation="Radarr",name="Radarr",sync_level="fullSync"} 1
		prowlarr_app_info{implementation="Sonarr",name="Sonarr",sync_level="fullSync"} 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_app_info"); err != nil {
		t.Error(err)
	}

	expected = `
		# HELP prowlarr_indexer_disabled Whether an indexer is temporarily disabled.
		# TYPE prowlarr_indexer_disabled gauge
		prowlarr_indexer_disabled{indexer="NZBgeek"} 1
	`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "prowlarr_indexer_disabled"); err != nil {
		t.Error(err)
	}
}

func TestSystemCollector_APIKeyHeader(t *testing.T) {
	var gotKey string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/system/status", func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-Api-Key")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version":"1.0","branch":"main","runtimeName":"test","runtimeVersion":""}`))
	})
	mux.HandleFunc("/api/v3/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	srv, client := newTestServer(mux)
	defer srv.Close()

	c := NewSystemCollector(client, 5*time.Second, slog.Default())
	// trigger collection
	testutil.CollectAndCount(c)

	if gotKey != "test-api-key" {
		t.Errorf("X-Api-Key = %q, want %q", gotKey, "test-api-key")
	}
}
