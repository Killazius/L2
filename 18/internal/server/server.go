package server

import (
	"18/internal/config"
	"18/internal/server/middleware"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler interface {
	CreateEvent(c *gin.Context)
	UpdateEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
	EventsForDay(c *gin.Context)
	EventsForWeek(c *gin.Context)
	EventsForMonth(c *gin.Context)
}

type Server struct {
	server  *http.Server
	cfg     config.HTTPConfig
	router  *gin.Engine
	handler Handler
}

func New(
	log *zap.SugaredLogger,
	cfg config.HTTPConfig,
	handler Handler,
) *Server {
	router := registerRoutes(log, handler)
	return &Server{
		server: &http.Server{
			Addr:         cfg.GetAddr(),
			ReadTimeout:  cfg.Timeout.Read,
			WriteTimeout: cfg.Timeout.Write,
			IdleTimeout:  cfg.Timeout.Idle,
			Handler:      router,
		},
		router:  router,
		cfg:     cfg,
		handler: handler,
	}

}

func registerRoutes(log *zap.SugaredLogger, h Handler) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger(log.Desugar()))

	api := r.Group("/api")
	{
		api.POST("/events", h.CreateEvent)
		api.PUT("/events", h.UpdateEvent)
		api.DELETE("/events", h.DeleteEvent)
		api.GET("/events/day", h.EventsForDay)
		api.GET("/events/week", h.EventsForWeek)
		api.GET("/events/month", h.EventsForMonth)
	}
	return r
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(fmt.Sprint("failed to start server: ", err))
	}
}
func (s *Server) Run() error {
	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout.Server)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
