package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"kayakaga-api/domain/mysql"

	"gorm.io/gorm"
)

type chatService struct {
	db        *gorm.DB
	agentURL  string
	httpClient *http.Client
}

func NewChatService(db *gorm.DB) UseCase {
	timeoutSecs, _ := strconv.Atoi(os.Getenv("AGENT_TIMEOUT"))
	if timeoutSecs == 0 {
		timeoutSecs = 60
	}

	return &chatService{
		db:   db,
		agentURL: os.Getenv("AGENT_API_URL"),
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSecs) * time.Second,
		},
	}
}

func (s *chatService) SendMessage(userID uint, token string, req *ChatRequest) (*ChatResponse, error) {
	// 1. Load user profile
	var user mysql.User
	err := s.db.Preload("Dependent").Preload("RiskProfile").First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load user profile: %w", err)
	}

	// 2. Load user accounts
	var accounts []mysql.Account
	err = s.db.Where("user_id = ?", userID).Preload("AccountType").Find(&accounts).Error
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
		UserToken:           token, // forward JWT token user
		UserContext:         userContext,
	}

	// 5. Call finai-agent-api
	agentResp, err := s.callAgent(&agentReq)
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

func (s *chatService) callAgent(req *AgentRequest) (*AgentResponse, error) {
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
