package rest

import (
	"encoding/json"
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger

	app *fiber.App
}

func New(log *zap.Logger) *Server {
	server := &Server{logger: log}

	server.app = fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	prometheus := fiberprometheus.New("auth")
	prometheus.RegisterAt(server.app, "/metrics")
	server.app.Use(prometheus.Middleware)

	server.app.Get("/healthz/liveness", server.liveness)
	server.app.Get("/healthz/readiness", server.readiness)

	return server
}

func (server *Server) Serve(port int) {
	server.logger.Info("HTTP server starts listening on", zap.Int("port", port))
	if err := server.app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		server.logger.Fatal("error starting HTTP server", zap.Error(err))
	}
}
