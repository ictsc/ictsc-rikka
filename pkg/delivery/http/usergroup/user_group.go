package usergroup

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupHandler struct {
	userGroupService service.UserGroupService
}

func NewUserGroupHandler(r *gin.RouterGroup, userGroupService service.UserGroupService) {
	handler := UserGroupHandler{
		userGroupService: userGroupService,
	}

	userGroup := r.Group("/user-groups")
	{
		authed := userGroup.Group("")
		authed.Use(middleware.AuthIsFullAccess())
		{
			authed.POST("", handler.Create)
		}
	}
}
