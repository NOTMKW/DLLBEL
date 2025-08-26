package models

import (
	"net"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type MT5Event struct {
	UserId    string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,porto3" json:"user_id,omitempty"`
	EventType string                 `protobuf:"bytes,2,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Symbol    string                 `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Volume    float64                `protobuf:"fixed64,4,opt,name=volume,proto3" json:"volume,omitempty"`
	Price     float64                `protobuf:"fixed64,5,opt,name=price,proto3" json:"price,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Data      []byte                 `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
}

type EnforcementMessage struct {
	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Action    string `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	Reason    string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	Severity  int32  `protobuf:"varint,4,opt,name=severity,proto3" json:"severity,omitempty"`
	Timestamp int64  `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
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
