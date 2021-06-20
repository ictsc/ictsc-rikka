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

type UserHandler struct {
	userController *controller.UserController
}

func NewUserHandler(r *gin.RouterGroup, userRepo repository.UserRepository, userService service.UserService) {
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
		}
	}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	req := &controller.CreateUserRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	res, err := h.userController.Create(req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *UserHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := h.userController.FindByID(id)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, "", nil, nil)
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}
