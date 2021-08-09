package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserHandler struct {
	userController *controller.UserController
}

func NewUserHandler(r *gin.RouterGroup, userRepo repository.UserRepository, userService *service.UserService) {
	handler := UserHandler{
		userController: controller.NewUserController(userService),
	}
	user := r.Group("/users")
	{
		user.POST("", handler.Create)

		authed := user.Group("")
		authed.Use(middleware.Auth(userRepo))
		{
			authed.GET("/:id", handler.FindByID)
			authed.PUT("/:id", handler.Update)
		}
	}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	req := &controller.CreateUserRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}

	res, err := h.userController.Create(req)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *UserHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := h.userController.FindByID(id)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	session := sessions.Default(ctx)
	signedInID, ok := session.Get("id").(string)
	if !ok {
		ctx.Error(error.NewUnauthorizedError("couldn't get id"))
		return
	}

	if id != signedInID {
		ctx.Error(error.NewForbiddenError("you can't update others information"))
		return
	}
	req := &controller.UpdateUserRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}

	res, err := h.userController.Update(id, req)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusAccepted, "", res, nil)
}
