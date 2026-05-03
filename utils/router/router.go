package router

import (
	"kayakaga-api/di"
	"kayakaga-api/modules/account"
	"kayakaga-api/modules/analytics"
	"kayakaga-api/modules/auth"
	"kayakaga-api/modules/goal"
	"kayakaga-api/modules/masters"
	"kayakaga-api/modules/simulate"
	"kayakaga-api/modules/transaction"
	"kayakaga-api/modules/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *di.Container) *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())

	auth.RegisterRoutes(r, container.Auth)

	protected := r.Group("/api/v1")
	protected.Use(AuthMiddleware())
	{
		user.RegisterRoutes(protected, container.User)
		account.RegisterRoutes(protected, container.Account)
		transaction.RegisterRoutes(protected, container.Transaction)
		goal.RegisterRoutes(protected, container.Goal)
		analytics.RegisterRoutes(protected, container.Analytics)
		simulate.RegisterRoutes(protected, container.Simulate)
		masters.RegisterRoutes(protected, container.Masters)
	}

	return r
}
