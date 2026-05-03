package chat

import "time"

// Request dari mobile
type ChatRequest struct {
	Message             string        `json:"message" binding:"required"`
	ConversationHistory []interface{} `json:"conversation_history"`
}

// Yang dikirim ke finai-agent-api
type AgentRequest struct {
	Message             string        `json:"message"`
	ConversationHistory []interface{} `json:"conversation_history"`
	UserToken           string        `json:"user_token"`
	UserContext         UserContext   `json:"user_context"`
}

type UserContext struct {
	Name          string    `json:"name"`
	City          string    `json:"city"`
	MonthlyIncome int64     `json:"monthly_income"`
	RiskProfile   string    `json:"risk_profile"`
	Accounts      []Account `json:"accounts"`
}

type Account struct {
	Name      string `json:"name"`
	Balance   int64  `json:"balance"`
	IsPrimary int8   `json:"is_primary"`
}

// Response dari finai-agent-api
type AgentResponse struct {
	Reply               string        `json:"reply"`
	ConversationHistory []interface{} `json:"conversation_history"`
	ToolsCalled         []ToolCall    `json:"tools_called"`
}

type ToolCall struct {
	Tool          string      `json:"tool"`
	Input         interface{} `json:"input"`
	ResultSummary string      `json:"result_summary"`
}

// Response ke mobile
type ChatResponse struct {
	Reply               string        `json:"reply"`
	ConversationHistory []interface{} `json:"conversation_history"`
	ToolsCalled         []ToolCall    `json:"tools_called"`
	Timestamp           time.Time     `json:"timestamp"`
}
