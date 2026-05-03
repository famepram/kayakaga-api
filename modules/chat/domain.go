package chat

type ChatUseCase interface {
	SendMessage(userID uint, token string, req *ChatRequest) (*ChatResponse, error)
}
