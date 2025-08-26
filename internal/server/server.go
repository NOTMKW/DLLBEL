package server

import (
	"context"
	"log"

	"github.com/NOTMKW/DLLBEL/internal/config"
	"github.com/NOTMKW/DLLBEL/internal/handlers"
	"github.com/NOTMKW/DLLBEL/internal/repository"
	"github.com/NOTMKW/DLLBEL/internal/routes"
	"github.com/NOTMKW/DLLBEL/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	app          *fiber.App
	config       *config.Config
	eventService *services.EventService
	repo         *repository.RedisRepository
}

func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	repo := repository.NewRedisRepository(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)

	ruleService := services.NewRuleService(repo)
	userService := services.NewUserService(repo)
	wsService := services.NewWebSocketService()

	eventService := services.NewWebSocketService(ruleService, userService, nil, wsService, cfg.EventBuffer)
	dllService := services.NewDLLService(eventService.GetEventChannel())

	eventService = services.NewWebSocketService(ruleService, userService, dllService, wsService, cfg.EventBuffer)

	wsHandler := handlers.NewWebSocketHandler(wsService)
	adminHandler := handlers.NewAdminHandler(ruleService, userService, dllService, wsService)
	dllHandler := handlers.NewDLLHandler(dllService)

	routes.SetupRoutes(app, wsHandler, adminHandler, dllHandler)

	eventService.Start(cfg.Workers)

	return &Server{
		app:          app,
		config:       cfg,
		eventService: eventService,
		repo:         repo,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(":" + s.config.Port)
}

func (s *Server) Shutdown() error {
	log.Println("Shutting down services...")

	s.eventService.Stop()
	s.repo.Close()

	return s.app.ShutdownWithContext(context.Background())
}
