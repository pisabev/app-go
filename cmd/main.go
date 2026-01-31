package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pisabev/app-go/common"
	"github.com/pisabev/app-go/http"
	"github.com/pisabev/app-go/service"

	"github.com/peterbourgon/ff/v3"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Load .env
	common.DotEnv()

	// Create context that gets canceled on certain signals
	ctx, cancelFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFn()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
	}
}

func run(ctx context.Context) error {
	fs := flag.NewFlagSet("app-go", flag.ExitOnError)

	config := common.Config{}

	fs.IntVar(&config.HttpPort, "http-port", 8080, "Http Server Port")
	fs.IntVar(&config.HttpPortDebug, "http-debug-port", 6060, "Http Debug Server Port")

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s \n", strings.Join(os.Args, " "))
		fs.PrintDefaults()
		return fmt.Errorf("flag set: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	srv := service.NewApp()

	h := http.NewServer(config.HttpPort, config.HttpPortDebug, srv)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return h.Serve(ctx) })
	g.Go(func() error { return h.ServeDebug(ctx) })

	return g.Wait()
}
