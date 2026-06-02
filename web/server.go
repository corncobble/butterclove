package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	readHeaderTimeout = 15 * time.Second
	shutdownTimeout   = 3 * time.Second
)

// StartServer will terminate program if http server could not be started.
// Note: This should be called in a separate goroutine due to blocking channel.
func StartServer(ctx context.Context, wg *sync.WaitGroup, host string, port int) {
	wg.Add(1)
	defer wg.Done()

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		ReadHeaderTimeout: readHeaderTimeout,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
		Handler:           handler,
	}

	// Start http server in separate goroutine.
	go func() {
		slog.InfoContext(ctx, "Start listen and serve.", "addr", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}

		slog.InfoContext(ctx, "Stop listen and serve.", "addr", server.Addr)
	}()

	// Block until context is canceled.
	<-ctx.Done()

	// Create an empty context and set the deadline.
	c, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Disable HTTP keep-alives.
	server.SetKeepAlivesEnabled(false)

	// Gracefully shut down the server.
	if err := server.Shutdown(c); err != nil {
		slog.WarnContext(ctx, "Unable to gracefully shutdown server", "err", err)
	}
}
