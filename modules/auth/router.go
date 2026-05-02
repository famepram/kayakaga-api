package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, uc UseCase) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", RegisterHandler(uc))
		auth.POST("/login", LoginHandler(uc))
		auth.POST("/refresh", RefreshHandler(uc))
		auth.POST("/logout", LogoutHandler(uc))
	}
}
