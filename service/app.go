package service

import (
	"context"
	"log/slog"
	"strings"
)

//go:generate moq -pkg mock -out ./mock/app.go . App
type App interface {
	Read(ctx context.Context, query string) (string, error)
	Something(ctx context.Context, str string) (string, error)
	SomethingId(ctx context.Context, str string, id string) (string, error)
}

type app struct {
	log *slog.Logger
}

func NewApp() App {
	return &app{
		log: slog.Default().With("component", "app"),
	}
}

func (t app) Read(ctx context.Context, query string) (string, error) {
	t.log.Info("Doing read", "input", query)
	return query, nil
}

func (t app) Something(ctx context.Context, str string) (string, error) {
	t.log.Info("Doing something", "input", str)
	return strings.ReplaceAll(str, " ", "-"), nil
}

func (t app) SomethingId(ctx context.Context, str string, id string) (string, error) {
	t.log.Info("Doing somethingId", "input", str)
	return id, nil
}
