package app

import (
	"18/internal/config"
	"18/internal/repository"
	"18/internal/server"
	"18/internal/server/handler"
	"18/internal/service"

	"go.uber.org/zap"
)

type App struct {
	log *zap.SugaredLogger
	api *server.Server
	cfg *config.Config
}

func New(log *zap.SugaredLogger, cfg *config.Config) *App {
	eventRepo := repository.New()
	eventService := service.New(eventRepo)
	h := handler.New(eventService)
	api := server.New(log, cfg.Server, h)
	return &App{
		log: log,
		api: api,
		cfg: cfg,
	}
}

func (a *App) Run() {
	defer func() {
		if r := recover(); r != nil {
			a.log.Errorw("application panicked and recovered", "panic", r)
			a.Stop()
		}
	}()
	a.log.Infow("start HTTP server", "addr", a.cfg.Server.GetAddr())
	a.api.MustRun()

}

func (a *App) Stop() {
	a.log.Info("closing HTTP server")
	if err := a.api.Close(); err != nil {
		a.log.Errorw("failed to stop HTTP server gracefully", "error", err)
	}
}
