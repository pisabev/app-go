package http

import (
	"encoding/json"
	"log/slog"
	httpserver "net/http"

	"github.com/asaskevich/govalidator"

	"github.com/pisabev/app-go/http/api"
)

func (s Server) read(w httpserver.ResponseWriter, r *httpserver.Request) {
	query := r.URL.Query().Get("query")

	res, err := s.srv.Read(r.Context(), query)
	if err != nil {
		s.Log.Warn("unable to handle request",
			slog.Any("payload", query),
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusInternalServerError, nil)
		return
	}

	s.Reply(w, httpserver.StatusOK, api.Response{
		Data: res,
	})
}

func (s Server) something(w httpserver.ResponseWriter, r *httpserver.Request) {
	var payload api.Request
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.Log.Warn("bad request",
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusBadRequest, nil)
		return
	}

	if v, err := govalidator.ValidateStruct(payload); !v {
		s.Log.Warn("bad request",
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusBadRequest, nil)
		return
	}

	res, err := s.srv.Something(r.Context(), payload.Field)
	if err != nil {
		s.Log.Warn("unable to handle request",
			slog.Any("payload", payload),
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusInternalServerError, nil)
		return
	}

	s.Reply(w, httpserver.StatusOK, api.Response{
		Data: res,
	})
}

func (s Server) somethingId(w httpserver.ResponseWriter, r *httpserver.Request) {
	id := r.PathValue("id")

	var payload api.Request
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.Log.Warn("bad request",
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusBadRequest, nil)
		return
	}

	if v, err := govalidator.ValidateStruct(payload); !v {
		s.Log.Warn("bad request",
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusBadRequest, nil)
		return
	}

	res, err := s.srv.SomethingId(r.Context(), payload.Field, id)
	if err != nil {
		s.Log.Warn("unable to handle request",
			slog.Any("payload", payload),
			slog.Any("error", err),
		)
		s.Reply(w, httpserver.StatusInternalServerError, nil)
		return
	}

	s.Reply(w, httpserver.StatusOK, api.Response{
		Data: res,
	})
}
