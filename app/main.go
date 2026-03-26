package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/rdmrcv/repartnersgo/app/api"
	"github.com/rdmrcv/repartnersgo/app/lifecycle"
)

const ShutdownDuration = time.Second * 5

// @BasePath /api
// @title Solver API
// @version 0.0.1

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	addr, has := os.LookupEnv("LISTEN_ADDR")
	if !has {
		addr = ":8080"
	}

	lifecycle.Run(
		context.Background(),
		logger,
		api.NewWorker(addr, logger),
		ShutdownDuration,
	)
}
