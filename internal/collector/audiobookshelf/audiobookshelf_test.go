package audiobookshelf

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
	"github.com/st0o0/tentacle/internal/collector"
)

func newTestClient(url string) *audiobookshelf.Client {
	return audiobookshelf.NewClient(url, "test-token", http.DefaultClient)
}

// --- SystemCollector ---

func TestSystemCollector_Up(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/api/backups", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"backups":[{"id":"b1","createdAt":1700000000000,"serverVersion":"2.17.0"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := newSystemCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_up", nil, 1)
	assertGauge(t, families, "audiobookshelf_system_info", map[string]string{"version": "2.17.0"}, 1)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "system"}, 1)
}

func TestSystemCollector_Down(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newSystemCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_up", nil, 0)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "system"}, 0)
}

func TestSystemCollector_Unreachable(t *testing.T) {
	c := newSystemCollector(newTestClient("http://127.0.0.1:1"), 1*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_up", nil, 0)
}

// --- LibrariesCollector ---

func TestLibrariesCollector(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/libraries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"libraries":[{"id":"lib1","name":"Audiobooks","mediaType":"book","lastUpdate":1700000000000},{"id":"lib2","name":"Podcasts","mediaType":"podcast","lastUpdate":1700500000000}]}`))
	})
	mux.HandleFunc("/api/libraries/lib1/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"totalItems":42,"totalSize":1073741824,"totalDuration":360000.5,"numAudioTracks":100,"totalAuthors":15,"totalGenres":8,"numMissing":2,"numInvalid":1}`))
	})
	mux.HandleFunc("/api/libraries/lib2/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"totalItems":10,"totalSize":524288000,"totalDuration":7200.0,"numAudioTracks":20,"totalAuthors":5,"totalGenres":3,"numMissing":0,"numInvalid":0}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := newLibrariesCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_libraries_total", nil, 2)
	assertGauge(t, families, "audiobookshelf_library_items_total", map[string]string{"library": "Audiobooks", "media_type": "book"}, 42)
	assertGauge(t, families, "audiobookshelf_library_size_bytes", map[string]string{"library": "Audiobooks", "media_type": "book"}, 1073741824)
	assertGauge(t, families, "audiobookshelf_library_items_total", map[string]string{"library": "Podcasts", "media_type": "podcast"}, 10)
	assertGauge(t, families, "audiobookshelf_library_size_bytes", map[string]string{"library": "Podcasts", "media_type": "podcast"}, 524288000)
	assertGauge(t, families, "audiobookshelf_library_last_update_timestamp_seconds", map[string]string{"library": "Audiobooks", "media_type": "book"}, 1700000000)
	assertGauge(t, families, "audiobookshelf_library_last_update_timestamp_seconds", map[string]string{"library": "Podcasts", "media_type": "podcast"}, 1700500000)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "libraries"}, 1)
}

func TestLibrariesCollector_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newLibrariesCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "libraries"}, 0)
}

// --- UsersCollector ---

func TestUsersCollector(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"u1","username":"alice","type":"admin","isActive":true,"lastSeen":1700000000000,"createdAt":1690000000000},{"id":"u2","username":"bob","type":"user","isActive":false,"lastSeen":1700500000000,"createdAt":1690000000000}]`))
	})
	mux.HandleFunc("/api/users/u1/listening-stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"totalTime":86400}`))
	})
	mux.HandleFunc("/api/users/u2/listening-stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"totalTime":3600}`))
	})
	mux.HandleFunc("/api/users/online", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"openSessions":[],"usersOnline":[{"id":"u1","username":"alice"}]}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := newUsersCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_users_total", nil, 2)
	assertGauge(t, families, "audiobookshelf_users_active_total", nil, 1)
	assertGauge(t, families, "audiobookshelf_users_online", nil, 1)
	assertGauge(t, families, "audiobookshelf_user_last_seen_timestamp_seconds", map[string]string{"user": "alice", "type": "admin"}, 1700000000)
	assertGauge(t, families, "audiobookshelf_user_last_seen_timestamp_seconds", map[string]string{"user": "bob", "type": "user"}, 1700500000)
	assertGauge(t, families, "audiobookshelf_user_listening_time_seconds", map[string]string{"user": "alice"}, 86400)
	assertGauge(t, families, "audiobookshelf_user_listening_time_seconds", map[string]string{"user": "bob"}, 3600)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "users"}, 1)
}

func TestUsersCollector_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newUsersCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "users"}, 0)
}

// --- SessionsCollector ---

func TestSessionsCollector(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sessions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total":57,"numPages":3,"page":0,"itemsPerPage":20,"sessions":[]}`))
	}))
	defer srv.Close()

	c := newSessionsCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_sessions_total", nil, 57)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "sessions"}, 1)
}

func TestSessionsCollector_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newSessionsCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "sessions"}, 0)
}

// --- BackupsCollector ---

func TestBackupsCollector(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/backups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"backups":[{"id":"b1","path":"/backups/b1.tar","datePretty":"2023-11-14","createdAt":1700000000000,"serverVersion":"2.7.0"},{"id":"b2","path":"/backups/b2.tar","datePretty":"2023-11-15","createdAt":1700100000000,"serverVersion":"2.7.0"},{"id":"b3","path":"/backups/b3.tar","datePretty":"2023-11-13","createdAt":1699900000000,"serverVersion":"2.7.0"}]}`))
	}))
	defer srv.Close()

	c := newBackupsCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_backups_total", nil, 3)
	// max CreatedAt: 1700100000000 ms -> 1700100000 seconds
	assertGauge(t, families, "audiobookshelf_backup_latest_timestamp_seconds", nil, 1700100000)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "backups"}, 1)
}

func TestBackupsCollector_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"backups":[]}`))
	}))
	defer srv.Close()

	c := newBackupsCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_backups_total", nil, 0)
	if findGauge(families, "audiobookshelf_backup_latest_timestamp_seconds", nil) != nil {
		t.Error("expected no latest_timestamp metric for empty backups")
	}
}

func TestBackupsCollector_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newBackupsCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())

	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "backups"}, 0)
}

// --- Auth header verification ---

func TestBearerTokenAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected 'Bearer test-token', got %q", auth)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total":0,"numPages":0,"page":0,"itemsPerPage":20,"sessions":[]}`))
	}))
	defer srv.Close()

	c := newSessionsCollector(newTestClient(srv.URL), 5*time.Second, slog.Default())
	families := collectSub(t, c)
	assertGauge(t, families, "audiobookshelf_scrape_success", map[string]string{"collector": "sessions"}, 1)
}

// --- helpers ---

func collectSub(t *testing.T, sc collector.SubCollector) map[string]*dto.MetricFamily {
	t.Helper()
	svc := collector.NewServiceCollector(namespace, slog.Default(), sc)
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(svc)
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
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
		if labelsMatch(m, labels) {
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

func labelsMatch(m *dto.Metric, expected map[string]string) bool {
	if len(expected) == 0 {
		return true
	}
	actual := make(map[string]string, len(m.GetLabel()))
	for _, lp := range m.GetLabel() {
		actual[lp.GetName()] = lp.GetValue()
	}
	for k, v := range expected {
		if actual[k] != v {
			return false
		}
	}
	return true
}
