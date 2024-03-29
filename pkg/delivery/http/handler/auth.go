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

type AuthHandler struct {
	authController *controller.AuthController
}

func NewAuthHandler(r *gin.RouterGroup, userRepo repository.UserRepository, authService *service.AuthService, userService *service.UserService) {
	handler := AuthHandler{
		authController: controller.NewAuthController(authService, userService),
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
	id, ok := session.Get("id").(string)
	if !ok {
		ctx.Error(error.NewUnauthorizedError("couldn't get id"))
		return
	}
	res, err := h.authController.Self(id)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	req := &controller.SignInRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return

	}

	session := sessions.Default(ctx)

	res, err := h.authController.SignIn(req)
	if err != nil {
		ctx.Error(error.NewUnauthorizedError(err.Error()))
		return
	}

	session.Set("id", res.User.ID.String())
	if err := session.Save(); err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AuthHandler) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})
	if err := session.Save(); err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}
	response.JSON(ctx, http.StatusOK, "", nil, nil)
}
