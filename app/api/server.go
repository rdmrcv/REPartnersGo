package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	gslog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
	"github.com/rdmrcv/repartnersgo/ui"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rdmrcv/repartnersgo/app/docs"
)

type worker struct {
	server *http.Server
	logger *slog.Logger
}

// Worker starts the server and listens for the cancellation signal
func (s *worker) Worker(ctx, shutdownCtx context.Context) {
	// since Shutdown makes ListenAndServe immediately return - we need way to wait
	// for Shutdown to stop its work and then return Worker function.
	fin := make(chan struct{})

	go func() {
		// If worker context stopped - initiate shutdown with special context
		<-ctx.Done()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("unclean http server shutdown", "error", err)
		}

		// Shutdown completed
		close(fin)
	}()

	if err := s.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("http server interrupted unexpectedly", "error", err)
	}

	select {
	case <-fin:
		// success shutdown
		s.logger.Info("http server stopped")
	case <-shutdownCtx.Done():
		// shutdownCtx canceled first — unclean finish
		s.logger.Error("unclean shutdown - http server cannot stop fast enough")
	}
}

// NewWorker create function that should be passed to the lifecycle.
func NewWorker(
	addr string,
	logger *slog.Logger,
) func(ctx, shutdownCtx context.Context) {
	router := gin.New()
	router.Use(
		gslog.SetLogger(),
		gin.Recovery(),
	)

	router.POST("/api/solve", Solve)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.StaticFS("/ui/", http.FS(ui.FS))

	return (&worker{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		logger: logger,
	}).Worker
}
