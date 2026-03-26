package main

import (
	"context"
	"fmt"
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

	port, has := os.LookupEnv("PORT")
	if !has {
		port = "8080"
	}

	lifecycle.Run(
		context.Background(),
		logger,
		api.NewWorker(fmt.Sprintf("0.0.0.0:%s", port), logger),
		ShutdownDuration,
	)
}
