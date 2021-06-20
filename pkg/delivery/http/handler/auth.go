package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type AuthHandler struct {
	authService    service.AuthService
	authController controller.AuthController
}

func NewAuthHandler(r *gin.RouterGroup, userRepo repository.UserRepository, authService service.AuthService, userService service.UserService) {
	handler := AuthHandler{
		authController: *controller.NewAuthController(&authService, &userService),
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

func (h *AuthHandler) Self(ctx *gin.Context) {
	session := sessions.Default(ctx)
	res, err := h.authController.Self(session)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, "", nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	req := &controller.SignInRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return

	}

	session := sessions.Default(ctx)

	res, err := h.authController.SignIn(req, session)
	if err != nil {
		response.JSON(ctx, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AuthHandler) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if err := h.authController.SignOut(session); err != nil {
		response.JSON(ctx, http.StatusInternalServerError, "", nil, nil)
		return
	}
	response.JSON(ctx, http.StatusOK, "", nil, nil)
}
