package services

import (
	"encoding/json"
	"sync"

	"github.com/NOTMKW/DLLBEL/internal/dto"
	"github.com/NOTMKW/DLLBEL/internal/models"
	"github.com/gofiber/websocket/v2"
)

type WebSocketService struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *WebSocketService) AddClient(id string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[id] = conn
}

func (s *WebSocketService) RemoveClient(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn, exists := s.clients[id]; exists {
		conn.Close()
		delete(s.clients, id)
	}
}

func (s *WebSocketService) SendMessage(id string, message interface{}) error {
	s.mu.RLock()
	conn, exists := s.clients[id]
	s.mu.RUnlock()
	if !exists {
		return nil
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	conn.WriteMessage(websocket.TextMessage, data)
	return nil
}

func (s *WebSocketService) BroadcastEvent(event *models.MT5Event) {
	message := dto.MT5EventMessage{
		Type: "mt5_event",
		Data: map[string]interface{}{
			"user_id":    event.UserID,
			"event_type": event.EventType,
			"symbol":     event.Symbol,
			"volume":     event.Volume,
			"price":      event.Price,
			"timestamp":  event.Timestamp,
		},
	}

	s.SendMessage(event.UserId, message)
}

func (s *WebSocketService) SendEnforcement(enforcement *models.EnforcementMessage) {
	message := dto.EnforcementMessage{
		Type: "enforcement",
		Data: map[string]interface{}{
			"user_id":   enforcement.UserId,
			"action":    enforcement.Action,
			"reason":    enforcement.Reason,
			"severity":  enforcement.Severity,
			"timestamp": enforcement.Timestamp,
		},
	}

	s.SendMessage(enforcement.UserId, message)
}

func (s *WebSocketService) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
