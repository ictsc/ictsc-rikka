package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(r *gin.RouterGroup, authService service.AuthService) {
	handler := AuthHandler{
		authService: authService,
	}
	auth := r.Group("auth")
	{
		auth.POST("/signin", handler.SignIn)

		authed := auth.Group("")
		authed.Use(middleware.Auth())
		{
			authed.DELETE("/signout", handler.SignOut)
		}
	}
}
