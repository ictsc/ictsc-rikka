package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupHandler struct {
	userGroupController *controller.UserGroupController
}

func NewUserGroupHandler(r *gin.RouterGroup, userRepo repository.UserRepository, userGroupService service.UserGroupService) {
	handler := UserGroupHandler{
		userGroupController: controller.NewUserGroupController(userGroupService),
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

func (h *UserGroupHandler) Create(ctx *gin.Context) {
	req := &controller.CreateUserGroupRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	res, err := h.userGroupController.Create(req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)

}
