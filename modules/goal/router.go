package goal

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	goals := r.Group("/goals")
	{
		goals.GET("", ListGoalsHandler(uc))
		goals.GET("/:id", GetGoalHandler(uc))
		goals.POST("", CreateGoalHandler(uc))
		goals.PUT("/:id", UpdateGoalHandler(uc))
		goals.DELETE("/:id", DeleteGoalHandler(uc))
		goals.PUT("/:id/contribution", UpdateContributionHandler(uc))
	}
}
