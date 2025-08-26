package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/NOTMKW/DLLBEL/internal/handlers"
)

func SetupRoutes(app *fiber.App, wsHandler *handlers.WebSocketHandler, adminHandler *handlers.AdminHandler, dllHandler *handlers.DLLHandler) {
	app.Get("/ws", websocket.New(wsHandler.HandleConnection))

	app.Post("/dll/connect", dllHandler.Connect)
	
	admin := app.Group("/admin")
	admin.Get("/rules", adminHandler.GetRules)
	admin.Post("/rules", adminHandler.CreateRule)
	admin.Put("/rules/:id", adminHandler.UpdateRule)
	admin.Delete("/rules/:id", adminHandler.DeleteRule)
	admin.Get("/users/:id/state", adminHandler.GetUserState)
	admin.Put("/users/:id/state", adminHandler.UpdateUserState)
	admin.Get("/connections", adminHandler.GetConnections)
	admin.Get("/metrics", adminHandler.GetMetrics)
	admin.Post("/enforce/:userid", adminHandler.ManualEnforce)
}