package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *ChatHandler) {
	chat := router.Group("/chat")
	{
		chat.POST("/message", handler.SendMessageHandler)
		chat.DELETE("/reset", handler.ResetChatHandler)
		chat.GET("/health", handler.HealthCheckHandler)
	}
}
