package masters

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	r.GET("/masters/account-types", GetAccountTypesHandler(uc))
	r.GET("/masters/categories", GetCategoriesHandler(uc))
	r.GET("/masters/goal-types", GetGoalTypesHandler(uc))
	r.GET("/masters/risk-profiles", GetRiskProfilesHandler(uc))
	r.GET("/masters/dependents", GetDependentsHandler(uc))
}
