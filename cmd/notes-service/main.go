package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/dngrudin/notes-service/internal/config"
	"github.com/dngrudin/notes-service/internal/service"
	"github.com/dngrudin/notes-service/internal/storage/memory"
	"github.com/dngrudin/notes-service/pkg/logger"
)

const (
	version = "0.0.1"
)

var (
	configPath   *string
	prettyLogger *bool
)

func init() {
	configPath = flag.String("config", "./config.yaml", "path to config file")
	prettyLogger = flag.Bool("pretty_log", true, "using pretty logger")
}

func main() {
	flag.Parse()

	cfg := config.MustLoad(*configPath)

	setupLogger(*prettyLogger)

	slog.Info("Starting Notes service", slog.String("version", version))

	storage := memory.New()

	service := service.New(cfg, storage)
	service.Run()
}

func setupLogger(isPrettyLogger bool) {
	var logHandler slog.Handler
	if isPrettyLogger {
		logHandler = logger.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	slog.SetDefault(slog.New(logHandler))
}
