package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(r *gin.RouterGroup, userRepo repository.UserRepository, authService service.AuthService) {
	handler := AuthHandler{
		authService: authService,
	}
	auth := r.Group("auth")
	{
		auth.POST("/signin", handler.SignIn)

		authed := auth.Group("")
		authed.Use(middleware.Auth(userRepo))
		{
			authed.GET("/self", handler.Self)
			authed.DELETE("/signout", handler.SignOut)
		}
	}
}
