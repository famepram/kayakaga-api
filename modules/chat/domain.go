package chat

type UseCase interface {
	SendMessage(userID uint, token string, req *ChatRequest) (*ChatResponse, error)
}
