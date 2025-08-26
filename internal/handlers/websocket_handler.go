package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/NOTMKW/DLLBEL/internal/services"
	"github.com/NOTMKW/DLLBEL/internal/dto"
)

type WebSocketHandler struct {
	wsService *services.WebSocketService
}

func NewWebSocketHandler(wsService *services.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		wsService: wsService,
	}
}

func (h *WebSocketHandler) HandleConnection(c *websocket.Conn) {
	userID := c.Query("user_id")
	if userID == "" {
		c.Close()
		return
	}

	h.wsService.AddClient(userID, c)
	defer h.wsService.RemoveClient(userID)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for user %s: %v", userID, err)
			break
		}

		h.processMessage(userID, message)
	}
}

func (h *WebSocketHandler) processMessage(userID string, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling WS message from %s: %v", userID, err)
		return
	}

	switch msg["type"] {
	case "ping":
		h.wsService.SendMessage(userID, &dto.WSMessage{
			Type: "pong",
			Data: nil,
		})
	case "subscribe":
	}
}