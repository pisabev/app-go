package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	httpserver "net/http"
	_ "net/http/pprof"
	"time"

	"app-go/service"
)

const (
	ApplicationJSON = "application/json"
	TextHTML        = "text/html"
)

type Server struct {
	Log       *slog.Logger
	Port      int
	PortDebug int
	srv       service.App
}

func NewServer(port, debug int, srv service.App) Server {
	return Server{
		Port:      port,
		PortDebug: debug,
		Log:       slog.Default().With("component", "http"),
		srv:       srv,
	}
}

func (s Server) setRoutes(r *httpserver.ServeMux) {
	r.Handle("GET /api/v1/read", http.HandlerFunc(s.read))
	r.Handle("POST /api/v1/something", http.HandlerFunc(s.something))
	r.Handle("POST /api/v1/something/{id}", http.HandlerFunc(s.somethingId))
}

func (s Server) Serve(ctx context.Context) error {
	r := httpserver.NewServeMux()
	s.setRoutes(r)

	return s._serve(ctx, &httpserver.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: s.RequestLogger(r),
	})
}

func (s Server) ServeDebug(ctx context.Context) error {
	return s._serve(ctx, &http.Server{
		Addr:    fmt.Sprintf(":%d", s.PortDebug),
		Handler: http.DefaultServeMux,
	})
}

func (s Server) _serve(ctx context.Context, srv *httpserver.Server) error {
	errCh := make(chan error, 1)

	go func() {
		s.Log.Info(fmt.Sprintf("Serving HTTP%s", srv.Addr))
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		s.Log.Info(fmt.Sprintf("Shutting down HTTP%s", srv.Addr))
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)

	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

func (s Server) Reply(w httpserver.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(code)
	if v != nil {
		err := json.NewEncoder(w).Encode(&v)
		if err != nil {
			s.Log.Error("Unable to encode", "error", err)
		}
	}
}

func (s Server) RequestLogger(next httpserver.Handler) httpserver.Handler {
	return httpserver.HandlerFunc(func(w httpserver.ResponseWriter, r *httpserver.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		s.Log.Info("executing",
			slog.String("url", r.URL.String()),
			slog.Int("port", s.Port),
			slog.String("ip", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.Float64("duration", float64(time.Since(startTime).Milliseconds())),
			slog.Any("headers", r.Header),
		)
	})
}
