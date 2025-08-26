package dto

type CreateRuleRequest struct {
	Name       string            `json:"name" validate:"required"`
	Conditions map[string]string `json:"conditions" validate:"required"`
	Actions    []string          `json:"actions" validate:"required"`
	Enabled    bool              `json:"enabled"`
	Priority   int               `json:"priority"`
}

type UpdateRuleRequest struct {
	Name       string            `json:"name"`
	Conditions map[string]string `json:"conditions"`
	Actions    []string          `json:"actions"`
	Enabled    *bool             `json:"enabled"`
	Priority   *int              `json:"priority"`
}

type UpdateUserStateRequest struct {
	Balance        *float64          `json:"balance"`
	Equity         *float64          `json:"equity"`
	OpenPositions  *int              `json:"open_positions"`
	DayVolume      *float64          `json:"day_volume"`
	RiskLevel      *string           `json:"risk_level"`
	ViolationCount *int              `json:"violation_count"`
	CustomData     map[string]string `json:"custom_data"`
}

type EnforceRequest struct {
	Action   string `json:"action" validate:"required"`
	Reason   string `json:"reason" validate:"required"`
	Severity int32  `json:"severity" validate:"min=1,max=5"`
}

type ConnectionInfo struct {
	ID       string `json:"id"`
	Active   bool   `json:"active"`
	LastPing int64  `json:"last_ping"`
}

type MetricsResponse struct {
	ActiveDLLConnections int   `json:"active_dll_connections"`
	WebSocketClients     int   `json:"websocket_clients"`
	UserStates           int   `json:"user_states"`
	EventBufferSize      int   `json:"event_buffer_size"`
	Timestamp            int64 `json:"timestamp"`
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
