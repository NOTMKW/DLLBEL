package models

import (
	"net"
	"sync"
	"time"
	"encoding/json"

)

type MT5Event struct {
	UserId    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Symbol    string    `json:"symbol"`
	Volume    float64   `json:"volume"`
	Price     float64   `json:"price"`
	Timestamp int64     `json:"timestamp"`
	Data      []byte    `json:"data,omitempty"`
}

type EnforcementMessage struct {
	UserId    string `json:"user_id"`
	Action    string `json:"action"`
	Reason    string `json:"reason"`
	Severity  int32  `json:"severity"`
	Timestamp int64  `json:"timestamp"`
}

func (e *MT5Event) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *MT5Event) Deserialize(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *EnforcementMessage) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EnforcementMessage) Deserialize(data []byte) error {
	return json.Unmarshal(data, e)
}

func NewMT5Event(userID, eventType, symbol string, volume, price float64) *MT5Event {
	return &MT5Event{
		UserId:    userID,
		EventType: eventType,
		Symbol:    symbol,
		Volume:    volume,
		Price:     price,
		Timestamp: time.Now().Unix(),
	}
}

func NewEnforcementMessage(userID, action, reason string, severity int32) *EnforcementMessage {
	return &EnforcementMessage{
		UserId:    userID,
		Action:    action,
		Reason:    reason,
		Severity:  severity,
		Timestamp: time.Now().Unix(),
	}
}

type Rule struct {
	ID         string            `json:"id" redis:"id"`
	Name       string            `json:"name" redis:"name"`
	Conditions map[string]string `json:"conditions" redis:"conditions"`
	Actions    []string          `json:"actions" redis:"actions"`
	Enabled    bool              `json:"enabled" redis:"enabled"`
	Priority   int               `json:"priority" redis:"priority"`
	CreatedAt  int64             `json:"created_at" redis:"created_at"`
	UpdatedAt  int64             `json:"updated_at" redis:"updated_at"`
}

type UserState struct {
	UserID         string            `json:"user_id" redis:"user_id"`
	Balance        float64           `json:"balance" redis:"balance"`
	Equity         float64           `json:"equity" redis:"equity"`
	OpenPositions  int               `json:"open_positions" redis:"open_positions"`
	DayVolume      float64           `json:"day_volume" redis:"day_volume"`
	LastActivity   int64             `json:"last_activity" redis:"last_activity"`
	RiskLevel      string            `json:"risk_level" redis:"risk_level"`
	ViolationCount int               `json:"violation_count" redis:"violation_count"`
	CustomData     map[string]string `json:"custom_data" redis:"custom_data"`
	Mu             sync.RWMutex      `json:"-" redis:"-"`
}

type DLLConnection struct {
	ID          string
	Conn        net.Conn
	IsActive    bool
	LastPing    int64
	EventChan   chan *MT5Event
	EnforceChan chan *EnforcementMessage
	Mu          sync.RWMutex
}
