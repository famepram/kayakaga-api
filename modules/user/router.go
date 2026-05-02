package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, uc UseCase) {
	users := r.Group("/users")
	{
		users.GET("/profile", GetProfileHandler(uc))
		users.PUT("/profile", UpdateProfileHandler(uc))
	}
}
