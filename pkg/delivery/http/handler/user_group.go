package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupHandler struct {
	userGroupController *controller.UserGroupController
}

func NewUserGroupHandler(r *gin.RouterGroup, userRepo repository.UserRepository, userService *service.UserService, userGroupService *service.UserGroupService, bastionService *service.BastionService) {
	handler := UserGroupHandler{
		userGroupController: controller.NewUserGroupController(userService, userGroupService, bastionService),
	}

	userGroup := r.Group("/usergroups")
	{
		authed := userGroup.Group("")
		authed.Use(middleware.Auth(userRepo))
		privileged := userGroup.Group("")
		privileged.Use(middleware.AuthIsFullAccess(userRepo))

		{
			authed.GET("", handler.ListParticipates)
			privileged.POST("", handler.Create)
		}
	}
}

func (h *UserGroupHandler) ListParticipates(ctx *gin.Context) {
	res, err := h.userGroupController.ListParticipates()
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *UserGroupHandler) Create(ctx *gin.Context) {
	req := &controller.CreateUserGroupRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}

	res, err := h.userGroupController.Create(req)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)

}
