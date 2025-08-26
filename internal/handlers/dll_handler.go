package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/NOTMKW/DLLBEL/internal/services")

type DLLHandler struct {
	dllService *services.DLLService
}

func NewDLLHandler(dllService *services.DLLService) *DLLHandler {
	return &DLLHandler{
		dllService: dllService,
	}
}

func (h *DLLHandler) Connect(c *fiber.Ctx) error {
	dllID := c.Query("dll_id")
	if dllID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "dll_id required"})
	}

	if err := h.dllService.StartListener(dllID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "connecting", "dll_id": dllID})
}
