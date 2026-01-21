package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	httpserver "net/http"
	"time"

	"app-go/service"
)

const (
	ApplicationJSON = "application/json"
	TextHTML        = "text/html"
)

type Server struct {
	Server *httpserver.Server
	Log    *slog.Logger
	Port   int
	srv    service.App
}

func NewServer(port int, srv service.App) Server {
	r := httpserver.NewServeMux()

	s := Server{
		Port: port,
		Log:  slog.Default().With("component", "http"),
		srv:  srv,
	}
	s.Server = &httpserver.Server{Handler: s.RequestLogger(r)}

	s.setRoutes(r)

	return s
}

func (s Server) setRoutes(r *httpserver.ServeMux) {

	r.Handle("GET /api/v1/read", http.HandlerFunc(s.read))
	r.Handle("POST /api/v1/something", http.HandlerFunc(s.something))
	r.Handle("POST /api/v1/something/{id}", http.HandlerFunc(s.somethingId))
}

func (s Server) Serve(ctx context.Context) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("http serve: %w", err)
	}

	go func() {
		<-ctx.Done()
		s.Log.Info("Shutting down HTTP...")
		_ = s.Server.Shutdown(context.Background())
	}()

	s.Log.Info(fmt.Sprintf("Serving HTTP on port %d", s.Port))

	if err := s.Server.Serve(ln); err != nil {
		return err
	}

	return nil
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
