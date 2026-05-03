package chat

import (
	"kayakaga-api/utils/helper"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	useCase ChatUseCase
}

func NewChatHandler(useCase ChatUseCase) *ChatHandler {
	return &ChatHandler{useCase: useCase}
}

// SendMessageHandler godoc
// @Summary Send chat message to AI
// @Description Send message to finai-agent-api for AI processing and get response
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ChatRequest true "Chat message with conversation history"
// @Success 200 {object} helper.Response{data=ChatResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 503 {object} helper.Response
// @Router /chat/message [post]
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

	resp, err := h.useCase.SendMessage(userID, token, &req)
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

// ResetChatHandler godoc
// @Summary Reset chat conversation
// @Description Reset/clear conversation history (client-side reset)
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Router /chat/reset [delete]
func (h *ChatHandler) ResetChatHandler(c *gin.Context) {
	// Client-side reset — server stateless
	// Hanya return acknowledgment
	helper.SuccessResponse(c, gin.H{
		"success": true,
		"message": "Conversation reset successfully",
	})
}

// HealthCheckHandler godoc
// @Summary Check chat service health
// @Description Check if kayakaga-api and finai-agent-api are healthy
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Router /chat/health [get]
func (h *ChatHandler) HealthCheckHandler(c *gin.Context) {
	// Check apakah finai-agent-api bisa dijangkau
	agentURL := os.Getenv("AGENT_API_URL")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(agentURL + "/health")

	agentStatus := "ok"
	if err != nil || resp.StatusCode != 200 {
		agentStatus = "unavailable"
	}
	if resp != nil {
		resp.Body.Close()
	}

	helper.SuccessResponse(c, gin.H{
		"kayakaga_api": "ok",
		"agent_api":    agentStatus,
	})
}
