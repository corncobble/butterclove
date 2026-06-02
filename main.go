package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/corncobble/butterclove/config"
	"github.com/corncobble/butterclove/web"
)

var (
	// Version is the version of this application, specified at compile time.
	Version = "0.0.0"
)

func main() {
	// Safely wait for all child routines to exit.
	var wg sync.WaitGroup

	// Top-level context will send a cancellation signal to child routines when its time to exit.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channel for our quit signal (SIGINT or SIGTERM).
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	// Read build info and set version.
	info, ok := debug.ReadBuildInfo()
	if ok {
		Version = info.Main.Version
	}

	slog.InfoContext(ctx, "Starting butterclove.", "version", Version)

	// Initialize configuration.
	config := config.NewConfig()
	if err := config.Read(); err != nil {
		slog.WarnContext(ctx, "Using default config.", "err", err)

		if err := config.Write(); err != nil {
			slog.ErrorContext(ctx, "Unable to write config.", "err", err)
		}
	}

	// Setup default server handler.
	web.SetupHandler(config.Channels)

	// Start goroutine for web server.
	go web.StartServer(ctx, &wg, config.Web.Host, config.Web.Port)

	// Block main thread until a quit signal is received.
	sig := <-quitCh

	// Stop incoming signals, close channel.
	signal.Stop(quitCh)
	close(quitCh)

	slog.InfoContext(ctx, "Stopping butterclove.", "signal", sig.String())

	// Cancel top-level context (signal to child routines to shutdown services and return).
	cancel()

	// Wait for all child routines to return (signal that they're done).
	wg.Wait()
}
