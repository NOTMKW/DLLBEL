package handlers

import (
	"time"

	"github.com/NOTMKW/DLLBEL/internal/dto"
	"github.com/NOTMKW/DLLBEL/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/NOTMKW/DLLBEL/internal/models"
)

type AdminHandler struct {
	ruleService *services.RuleService
	wsService   *services.WebSocketService
	dllService  *services.DLLService
	userService *services.UserService
}

func NewAdminHandler(ruleService *services.RuleService, wsService *services.WebSocketService, dllService *services.DLLService, userService *services.UserService) *AdminHandler {
	return &AdminHandler{
		ruleService: ruleService,
		wsService:   wsService,
		dllService:  dllService,
		userService: userService,
	}
}

func (h *AdminHandler) GetRules(c *fiber.Ctx) error {
	rules, err := h.ruleService.GetAllRules()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch rules"})
	}
	return c.JSON(rules)
}

func (h *AdminHandler) CreateRule(c *fiber.Ctx) error {
	var req dto.CreateRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	rule, err := h.ruleService.CreateRule(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rule)
}

func (h *AdminHandler) UpdateRule(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	rule, err := h.ruleService.UpdateRule(id, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rule)
}

func (h *AdminHandler) DeleteRule(c *fiber.Ctx) error {
	id := c.Params("id")
	
	h.ruleService.DeleteRule(id)
	return c.JSON(fiber.Map{"message": "Rule deleted"})
}

func (h *AdminHandler) GetUserState(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	state := h.userService.GetUserState(userID)
	if state == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User state not found"})
	}
	return c.JSON(state)
}

func (h *AdminHandler) UpdateUserState(c *fiber.Ctx) error {
	userID := c.Params("id")
	
	var req dto.UpdateUserStateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	state := h.userService.UpdateUserState(userID, &req)
	return c.JSON(state)
}

func (h *AdminHandler) GetConnections(c *fiber.Ctx) error {
	conns := h.dllService.GetConnections()
	return c.JSON(conns)
}

func (h *AdminHandler) GetMetrics(c *fiber.Ctx) error {
	metrics := &dto.MetricsResponse{
		ActiveDLLConnections: h.dllService.GetActiveConnectionCount(),
		WebSocketClients:     h.wsService.GetClientCount(),
		UserStates:          h.userService.GetUserCount(),
		EventBufferSize:     0, // Will be set by the calling service
		Timestamp:           time.Now().Unix(),
	}

	return c.JSON(metrics)
}

func (h *AdminHandler) ManualEnforce(c *fiber.Ctx) error {
	userID := c.Params("userid")
	
	var req dto.EnforceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	enforcement := &models.EnforcementMessage{
		UserId:    userID,
		Action:    req.Action,
		Reason:    req.Reason,
		Severity:  req.Severity,
		Timestamp: time.Now().Unix(),
	}

	h.dllService.SendEnforcement(enforcement)
	h.wsService.SendEnforcement(enforcement)

	return c.JSON(fiber.Map{"status": "sent", "enforcement": enforcement})
}