package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServe(ctx context.Context, addr string, reg *prometheus.Registry, timeout time.Duration, version string, logger *slog.Logger) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Timeout:       timeout,
		ErrorHandling: promhttp.ContinueOnError,
	})

	mux := http.NewServeMux()
	mux.Handle("/metrics", requestLogger(metricsHandler, logger))
	mux.HandleFunc("/healthz", healthHandler(version))

	srv := &http.Server{Handler: mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("metrics server shutdown", "err", err)
		}
	}()

	logger.Info("metrics server listening", "addr", ln.Addr().String())
	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func requestLogger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		logger.Info("scrape complete", "status", sw.status, "duration", time.Since(start).Seconds(), "client", r.RemoteAddr)
	})
}

func healthHandler(version string) http.HandlerFunc {
	start := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"version": version,
			"uptime":  time.Since(start).String(),
		})
	}
}
