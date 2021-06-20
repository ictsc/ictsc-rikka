package usergroup

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupHandler struct {
	userRepo         repository.UserRepository
	userGroupService service.UserGroupService
}

func NewUserGroupHandler(r *gin.RouterGroup, userRepo repository.UserRepository, userGroupService service.UserGroupService) {
	handler := UserGroupHandler{
		userRepo:         userRepo,
		userGroupService: userGroupService,
	}

	userGroup := r.Group("/user-groups")
	{
		authed := userGroup.Group("")
		authed.Use(middleware.AuthIsFullAccess(userRepo))
		{
			authed.POST("", handler.Create)
		}
	}
}
