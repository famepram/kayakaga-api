package simulate

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	r.GET("/simulate/investment", SimulateInvestmentHandler(uc))
}
