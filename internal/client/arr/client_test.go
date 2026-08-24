package arr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type recorded struct {
	method string
	path   string
	header http.Header
	query  string
}

func newTestServer(t *testing.T, statusCode int, body any) (*httptest.Server, *recorded) {
	t.Helper()
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.path = r.URL.Path
		rec.header = r.Header
		rec.query = r.URL.RawQuery

		w.WriteHeader(statusCode)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, rec
}

func TestAuthHeader(t *testing.T) {
	srv, rec := newTestServer(t, http.StatusOK, &SystemStatus{Version: "4.0"})
	c := NewClient(srv.URL, "test-key-123", srv.Client())

	_, _ = c.GetSystemStatus(context.Background())

	got := rec.header.Get("X-Api-Key")
	if got != "test-key-123" {
		t.Fatalf("X-Api-Key = %q, want %q", got, "test-key-123")
	}
}

func TestURLConstruction(t *testing.T) {
	srv, rec := newTestServer(t, http.StatusOK, &SystemStatus{})

	c := NewClient(srv.URL+"/", "key", srv.Client())
	_, _ = c.GetSystemStatus(context.Background())

	want := "/api/v3/system/status"
	if rec.path != want {
		t.Fatalf("path = %q, want %q", rec.path, want)
	}
}

func TestErrorHandling(t *testing.T) {
	srv, _ := newTestServer(t, http.StatusUnauthorized, nil)
	c := NewClient(srv.URL, "bad-key", srv.Client())

	_, err := c.GetSystemStatus(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	want := "unexpected status 401"
	if got := err.Error(); !contains(got, want) {
		t.Fatalf("error = %q, want substring %q", got, want)
	}
}

func TestGetSystemStatus(t *testing.T) {
	payload := SystemStatus{
		Version:   "4.7.2",
		Branch:    "main",
		StartTime: "2024-01-01T00:00:00Z",
	}
	srv, _ := newTestServer(t, http.StatusOK, &payload)
	c := NewClient(srv.URL, "key", srv.Client())

	status, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Version != "4.7.2" {
		t.Fatalf("Version = %q, want %q", status.Version, "4.7.2")
	}
	if status.Branch != "main" {
		t.Fatalf("Branch = %q, want %q", status.Branch, "main")
	}
}

func TestGetHealth(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		payload := []HealthCheck{
			{Source: "IndexerCheck", Type: "warning", Message: "no indexers"},
		}
		srv, _ := newTestServer(t, http.StatusOK, payload)
		c := NewClient(srv.URL, "key", srv.Client())

		checks, err := c.GetHealth(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(checks) != 1 {
			t.Fatalf("got %d checks, want 1", len(checks))
		}
		if checks[0].Source != "IndexerCheck" {
			t.Fatalf("Source = %q, want %q", checks[0].Source, "IndexerCheck")
		}
	})

	t.Run("empty array", func(t *testing.T) {
		srv, _ := newTestServer(t, http.StatusOK, []HealthCheck{})
		c := NewClient(srv.URL, "key", srv.Client())

		checks, err := c.GetHealth(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(checks) != 0 {
			t.Fatalf("got %d checks, want 0", len(checks))
		}
	})
}

func TestGetQueueParams(t *testing.T) {
	srv, rec := newTestServer(t, http.StatusOK, &QueueResponse{TotalRecords: 5})
	c := NewClient(srv.URL, "key", srv.Client())

	resp, err := c.GetQueue(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.path != "/api/v3/queue" {
		t.Fatalf("path = %q, want /api/v3/queue", rec.path)
	}

	q := rec.query
	if !contains(q, "page=1") {
		t.Fatalf("query %q missing page=1", q)
	}
	if !contains(q, "pageSize=250") {
		t.Fatalf("query %q missing pageSize=250", q)
	}
	if resp.TotalRecords != 5 {
		t.Fatalf("TotalRecords = %d, want 5", resp.TotalRecords)
	}
}

func TestGetRootFolders(t *testing.T) {
	payload := []RootFolder{
		{Path: "/data/movies", FreeSpace: 1000000, TotalSpace: 5000000},
		{Path: "/data/tv", FreeSpace: 2000000, TotalSpace: 8000000},
	}
	srv, _ := newTestServer(t, http.StatusOK, payload)
	c := NewClient(srv.URL, "key", srv.Client())

	folders, err := c.GetRootFolders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(folders) != 2 {
		t.Fatalf("got %d folders, want 2", len(folders))
	}
	if folders[0].FreeSpace != 1000000 {
		t.Fatalf("FreeSpace = %d, want 1000000", folders[0].FreeSpace)
	}
	if folders[1].TotalSpace != 8000000 {
		t.Fatalf("TotalSpace = %d, want 8000000", folders[1].TotalSpace)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstr(s, substr)
}

func searchSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
