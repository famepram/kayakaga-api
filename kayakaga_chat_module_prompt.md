# Prompt 2 — Chat Module di kayakaga-api (Golang)

Tambahkan module baru `chat` di kayakaga-api.
Module ini menjadi gateway antara mobile app dan finai-agent-api.
Ikuti struktur module yang sudah ada.

---

## Struktur baru

```
modules/
└── chat/
    ├── domain.go      # Interface ChatUseCase
    ├── entities.go    # Request/response structs
    ├── handler.go     # HTTP handlers
    ├── service.go     # Business logic + call ke finai-agent-api
    └── router.go      # Route definitions
```

---

## .env additions

```
AGENT_API_URL=http://localhost:8000    # finai-agent-api URL
AGENT_TIMEOUT=60                       # timeout dalam detik (AI bisa lambat)
```

---

## modules/chat/entities.go

```go
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
    Tool          string `json:"tool"`
    Input         interface{} `json:"input"`
    ResultSummary string `json:"result_summary"`
}

// Response ke mobile
type ChatResponse struct {
    Reply               string        `json:"reply"`
    ConversationHistory []interface{} `json:"conversation_history"`
    ToolsCalled         []ToolCall    `json:"tools_called"`
    Timestamp           time.Time     `json:"timestamp"`
}
```

---

## modules/chat/domain.go

```go
package chat

type ChatUseCase interface {
    SendMessage(userID uint, token string, req ChatRequest) (*ChatResponse, error)
}
```

---

## modules/chat/service.go

```go
package chat

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "time"

    "kayakaga-api/domain/mysql/providers"
)

type chatService struct {
    userProvider    providers.UserProvider
    accountProvider providers.AccountProvider
    agentURL        string
    httpClient      *http.Client
}

func NewChatService(
    userProvider providers.UserProvider,
    accountProvider providers.AccountProvider,
) ChatUseCase {
    timeoutSecs, _ := strconv.Atoi(os.Getenv("AGENT_TIMEOUT"))
    if timeoutSecs == 0 {
        timeoutSecs = 60
    }

    return &chatService{
        userProvider:    userProvider,
        accountProvider: accountProvider,
        agentURL:        os.Getenv("AGENT_API_URL"),
        httpClient: &http.Client{
            Timeout: time.Duration(timeoutSecs) * time.Second,
        },
    }
}

func (s *chatService) SendMessage(userID uint, token string, req ChatRequest) (*ChatResponse, error) {
    // 1. Load user profile
    user, err := s.userProvider.GetByID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to load user profile: %w", err)
    }

    // 2. Load user accounts
    accounts, err := s.accountProvider.ListByUserID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to load accounts: %w", err)
    }

    // 3. Build user context untuk agent
    accountList := make([]Account, len(accounts))
    for i, acc := range accounts {
        accountList[i] = Account{
            Name:      acc.Name,
            Balance:   acc.Balance,
            IsPrimary: acc.IsPrimary,
        }
    }

    // Get risk profile code
    riskProfileCode := "undecided"
    if user.RiskProfile.Code != "" {
        riskProfileCode = user.RiskProfile.Code
    }

    userContext := UserContext{
        Name:          user.Name,
        City:          user.City,
        MonthlyIncome: user.MonthlyIncome,
        RiskProfile:   riskProfileCode,
        Accounts:      accountList,
    }

    // 4. Build request ke finai-agent-api
    agentReq := AgentRequest{
        Message:             req.Message,
        ConversationHistory: req.ConversationHistory,
        UserToken:           token,   // forward JWT token user
        UserContext:         userContext,
    }

    // 5. Call finai-agent-api
    agentResp, err := s.callAgent(agentReq)
    if err != nil {
        return nil, fmt.Errorf("agent service error: %w", err)
    }

    // 6. Return response ke mobile
    return &ChatResponse{
        Reply:               agentResp.Reply,
        ConversationHistory: agentResp.ConversationHistory,
        ToolsCalled:         agentResp.ToolsCalled,
        Timestamp:           time.Now().UTC(),
    }, nil
}

func (s *chatService) callAgent(req AgentRequest) (*AgentResponse, error) {
    body, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    resp, err := s.httpClient.Post(
        fmt.Sprintf("%s/agent/chat", s.agentURL),
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return nil, fmt.Errorf("cannot reach agent service: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("agent returned status %d", resp.StatusCode)
    }

    var agentResp AgentResponse
    if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
        return nil, fmt.Errorf("failed to decode agent response: %w", err)
    }

    return &agentResp, nil
}
```

---

## modules/chat/handler.go

```go
package chat

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "kayakaga-api/utils/helper"
)

type ChatHandler struct {
    useCase ChatUseCase
}

func NewChatHandler(useCase ChatUseCase) *ChatHandler {
    return &ChatHandler{useCase: useCase}
}

func (h *ChatHandler) SendMessageHandler(c *gin.Context) {
    userID := c.GetUint("user_id")
    if userID == 0 {
        helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
        return
    }

    // Extract raw JWT token dari header
    authHeader := c.GetHeader("Authorization")
    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == "" {
        helper.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "No token provided")
        return
    }

    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
        return
    }

    if strings.TrimSpace(req.Message) == "" {
        helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "Message cannot be empty")
        return
    }

    // Initialize empty history jika nil
    if req.ConversationHistory == nil {
        req.ConversationHistory = []interface{}{}
    }

    resp, err := h.useCase.SendMessage(userID, token, req)
    if err != nil {
        // Differentiate error types
        errMsg := err.Error()
        if strings.Contains(errMsg, "cannot reach agent") {
            helper.ErrorResponse(c, http.StatusServiceUnavailable,
                "AGENT_UNAVAILABLE",
                "AI service is currently unavailable. Please try again.")
            return
        }
        helper.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", errMsg)
        return
    }

    helper.SuccessResponse(c, resp)
}

func (h *ChatHandler) ResetChatHandler(c *gin.Context) {
    // Client-side reset — server stateless
    // Hanya return acknowledgment
    helper.SuccessResponse(c, gin.H{
        "success": true,
        "message": "Conversation reset successfully",
    })
}

func (h *ChatHandler) HealthCheckHandler(c *gin.Context) {
    // Check apakah finai-agent-api bisa dijangkau
    import (
        "net/http"
        "os"
        "time"
    )

    agentURL := os.Getenv("AGENT_API_URL")
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(agentURL + "/health")

    agentStatus := "ok"
    if err != nil || resp.StatusCode != 200 {
        agentStatus = "unavailable"
    }

    helper.SuccessResponse(c, gin.H{
        "kayakaga_api": "ok",
        "agent_api":    agentStatus,
    })
}
```

---

## modules/chat/router.go

```go
package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(
    router *gin.RouterGroup,
    handler *ChatHandler,
) {
    chat := router.Group("/chat")
    {
        chat.POST("/message", handler.SendMessageHandler)
        chat.DELETE("/reset", handler.ResetChatHandler)
        chat.GET("/health", handler.HealthCheckHandler)
    }
}
```

---

## Endpoints yang ditambahkan

```
POST   /api/v1/chat/message    ← mobile kirim pesan AI
DELETE /api/v1/chat/reset      ← reset conversation
GET    /api/v1/chat/health     ← cek status agent service
```

---

## Register di utils/router/router.go

Tambahkan chat module mengikuti pattern yang sudah ada:

```go
// Import chat module
chatHandler := chat.NewChatHandler(
    chat.NewChatService(userProvider, accountProvider),
)

// Register routes (protected — butuh JWT)
protected := v1.Group("")
protected.Use(middleware.JWTAuth())
{
    // ... existing routes ...
    chat.RegisterRoutes(protected, chatHandler)
}
```

---

## Register di di/domain.go

Tambahkan provider untuk chat module mengikuti pattern DI yang sudah ada.

---

## Error handling yang perlu dicover

```
1. finai-agent-api down          → 503 AGENT_UNAVAILABLE
2. finai-agent-api timeout       → 504 GATEWAY_TIMEOUT  
3. User profile tidak ditemukan  → 404 USER_NOT_FOUND
4. Message kosong                → 400 INVALID_INPUT
5. Token tidak valid             → 401 UNAUTHORIZED
```

---

## Test setelah build

```bash
TOKEN="paste_jwt_token_disini"

# 1. Health check
curl http://localhost:8080/api/v1/chat/health \
  -H "Authorization: Bearer $TOKEN"

# Expected:
# {"success":true,"data":{"kayakaga_api":"ok","agent_api":"ok"}}

# 2. Send message
curl -X POST http://localhost:8080/api/v1/chat/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Cek saldo semua akun",
    "conversation_history": []
  }'

# Expected:
# {
#   "success": true,
#   "data": {
#     "reply": "Total saldo kamu Rp 23.4jt...",
#     "conversation_history": [...],
#     "tools_called": [...]
#   }
# }

# 3. Multi-turn conversation
curl -X POST http://localhost:8080/api/v1/chat/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "GoPay cukup sampai akhir bulan?",
    "conversation_history": [dari response sebelumnya]
  }'

# 4. Reset
curl -X DELETE http://localhost:8080/api/v1/chat/reset \
  -H "Authorization: Bearer $TOKEN"
```

---

## Important notes

1. Timeout untuk chat endpoint harus lebih panjang dari endpoint lain
   (AI response bisa 5-15 detik) — set 60 detik di HTTP client

2. Jangan cache response AI — setiap request harus fresh

3. conversation_history di-pass as-is dari mobile ke agent dan balik
   kayakaga-api tidak perlu parse isinya

4. user_token yang di-forward ke agent adalah token yang sama
   yang dipakai mobile — agent pakai ini untuk call balik ke kayakaga-api

5. Pastikan AGENT_API_URL tidak exposed ke public
   Hanya accessible dari internal network / localhost
