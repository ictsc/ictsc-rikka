package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(r *gin.RouterGroup, userService service.UserService) {
	handler := UserHandler{
		userService: userService,
	}
	user := r.Group("/users")
	{
		user.POST("", handler.Create)

		authed := user.Group("")
		authed.Use(middleware.Auth())
		{
			authed.GET("/:id", handler.FindByID)
		}
	}
}
