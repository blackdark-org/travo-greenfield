// Command travo-greenfield is the on-router daemon.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BlackDark/travo-greenfield/internal/api"
	"github.com/BlackDark/travo-greenfield/internal/app"
	"github.com/BlackDark/travo-greenfield/internal/version"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8081", "listen address")
	useMock := flag.Bool("mock", false, "run against the mock device (no router needed)")
	flag.Parse()

	deps, err := app.New(*useMock)
	if err != nil {
		slog.Error("startup", "err", err)
		os.Exit(1)
	}
	slog.Info("adapters", "mock", deps.Mock, "packages", deps.Packages.Name())

	h, _ := api.New(version.Version)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	slog.Info("listening", "addr", *addr, "version", version.Version)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server", "err", err)
		os.Exit(1)
	}
}
