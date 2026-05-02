package transaction

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	txns := r.Group("/transactions")
	{
		txns.GET("", ListTransactionsHandler(uc))
		txns.POST("", CreateTransactionHandler(uc))
		txns.PUT("/:id", UpdateTransactionHandler(uc))
		txns.DELETE("/:id", DeleteTransactionHandler(uc))
	}
}
