package lifecycle

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run wraps the main application run into the lifecycle where we can organize
// graceful shutdown.
//
// Worker should handle their own errors themselves. If the worker cannot handle
// an error — it should properly log the error and returns.
func Run(
	ctx context.Context,
	logger *slog.Logger,
	worker func(context.Context, context.Context),
	terminationTimeout time.Duration,
) {
	if worker == nil {
		os.Exit(1)
	}

	// Create work ctx to allow graceful stop.
	workCtx, cancel := context.WithCancel(ctx)

	go func() {
		// if the worker stopped, we cancel work context and initiate shutdown
		defer func() { cancel() }()

		worker(workCtx, ctx)
	}()

	// listen for shutdown signals and handle unexpected worker stops (if worker
	// stopped themselves without our signal) and background context cancellation
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		// if signal received, initiate shutdown and wait worker to be stopped
		cancel()

		logger.Info("shutting down...")

		// clearly separate worker stop, timeout and background context interruption
		select {
		case <-workCtx.Done():
			logger.Info("worker stopped")
		case <-time.After(terminationTimeout):
			logger.Warn("unclean shutdown due to timeout")
		case <-ctx.Done():
			logger.Warn("unclean shutdown")
		}
	case <-workCtx.Done():
		logger.Warn("worker unexpectedly terminated")
	case <-ctx.Done():
		logger.Warn("unclean shutdown")
	}
}
