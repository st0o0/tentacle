package seerr

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/st0o0/tentacle/internal/client/seerr"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(seerr.Status{
			Version:   "2.3.0",
			CommitTag: "abc1234",
		})
	})

	mux.HandleFunc("/api/v1/request/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(seerr.RequestCount{
			Total:     42,
			Movie:     25,
			TV:        17,
			Pending:   5,
			Approved:  10,
			Available: 20,
			Declined:  7,
		})
	})

	mux.HandleFunc("/api/v1/issue/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(seerr.IssueCount{
			Total:    10,
			Open:     3,
			Resolved: 7,
		})
	})

	mux.HandleFunc("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(seerr.UsersResponse{
			PageInfo: seerr.PageInfo{
				Pages:   1,
				Results: 15,
			},
		})
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func newClient(srv *httptest.Server) *seerr.Client {
	return seerr.NewClient(srv.URL, "test-api-key", srv.Client())
}

func TestSystemCollector(t *testing.T) {
	srv := newTestServer(t)
	c := NewSystemCollector(newClient(srv), 5*time.Second, slog.Default())

	expected := `
		# HELP seerr_up Whether Seerr is reachable.
		# TYPE seerr_up gauge
		seerr_up 1
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_up"); err != nil {
		t.Fatal(err)
	}

	expected = `
		# HELP seerr_system_info Seerr system information.
		# TYPE seerr_system_info gauge
		seerr_system_info{commit="abc1234",version="2.3.0"} 1
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_system_info"); err != nil {
		t.Fatal(err)
	}
}

func TestRequestsCollector(t *testing.T) {
	srv := newTestServer(t)
	c := NewRequestsCollector(newClient(srv), 5*time.Second, slog.Default())

	expected := `
		# HELP seerr_requests_total Total number of Seerr requests.
		# TYPE seerr_requests_total gauge
		seerr_requests_total 42
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_requests_total"); err != nil {
		t.Fatal(err)
	}

	expected = `
		# HELP seerr_requests_by_type_total Number of Seerr requests by media type.
		# TYPE seerr_requests_by_type_total gauge
		seerr_requests_by_type_total{media_type="movie"} 25
		seerr_requests_by_type_total{media_type="tv"} 17
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_requests_by_type_total"); err != nil {
		t.Fatal(err)
	}

	expected = `
		# HELP seerr_requests_by_status_total Number of Seerr requests by status.
		# TYPE seerr_requests_by_status_total gauge
		seerr_requests_by_status_total{status="pending"} 5
		seerr_requests_by_status_total{status="approved"} 10
		seerr_requests_by_status_total{status="available"} 20
		seerr_requests_by_status_total{status="declined"} 7
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_requests_by_status_total"); err != nil {
		t.Fatal(err)
	}
}

func TestUsersCollector(t *testing.T) {
	srv := newTestServer(t)
	c := NewUsersCollector(newClient(srv), 5*time.Second, slog.Default())

	expected := `
		# HELP seerr_users_total Total number of Seerr users.
		# TYPE seerr_users_total gauge
		seerr_users_total 15
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_users_total"); err != nil {
		t.Fatal(err)
	}
}

func TestIssuesCollector(t *testing.T) {
	srv := newTestServer(t)
	c := NewIssuesCollector(newClient(srv), 5*time.Second, slog.Default())

	expected := `
		# HELP seerr_issues_total Total number of reported issues.
		# TYPE seerr_issues_total gauge
		seerr_issues_total 10
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_issues_total"); err != nil {
		t.Fatal(err)
	}

	expected = `
		# HELP seerr_issues_by_status_total Number of issues by status.
		# TYPE seerr_issues_by_status_total gauge
		seerr_issues_by_status_total{status="open"} 3
		seerr_issues_by_status_total{status="resolved"} 7
	`
	if err := testutil.CollectAndCompare(c, readExpected(expected), "seerr_issues_by_status_total"); err != nil {
		t.Fatal(err)
	}
}

func readExpected(s string) *strings.Reader {
	return strings.NewReader(s)
}
