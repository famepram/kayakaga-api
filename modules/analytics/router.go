package analytics

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	analytics := r.Group("/analytics")
	{
		analytics.GET("/budget", GetBudgetHandler(uc))
		analytics.GET("/compare", GetCompareHandler(uc))
		analytics.GET("/anomalies", GetAnomaliesHandler(uc))
		analytics.GET("/recurring", GetRecurringHandler(uc))
		analytics.GET("/savings-suggestion", GetSavingsSuggestionHandler(uc))
		analytics.GET("/goal-recommendation", GetGoalRecommendationHandler(uc))
	}
}
