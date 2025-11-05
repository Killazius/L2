package main

import (
	"18/internal/app"
	"18/internal/config"
	"18/internal/logger"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log, err := logger.LoadFromConfig(cfg.Logger.Path)
	if err != nil {
		if errors.Is(err, logger.ErrDefaultLogger) {
			log.Warnw("using default logger because config file not found",
				"config_path", cfg.Logger.Path)
		} else {
			panic(fmt.Sprintf("logger load error: %s", err))
		}
	}
	application := app.New(log, cfg)
	go application.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	application.Stop()
	log.Infow("application stopped")
}
