package account

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	accounts := r.Group("/accounts")
	{
		accounts.GET("", ListAccountsHandler(uc))
		accounts.GET("/balances", GetBalancesHandler(uc))
		accounts.POST("", CreateAccountHandler(uc))
		accounts.PUT("/:id", UpdateAccountHandler(uc))
		accounts.DELETE("/:id", DeleteAccountHandler(uc))
	}
}
