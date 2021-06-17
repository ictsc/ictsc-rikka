package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
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

type SignInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	User    *entity.User `json:"user,omitempty"`
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	req := SignInRequest{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, SignInResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return

	}
	user, err := h.authService.SignIn(req.Name, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, SignInResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	session := sessions.Default(ctx)
	session.Set("id", user.ID.String())
	session.Set("isFullAccess", user.UserGroup.IsFullAccess)
	session.Set("isReadOnly", user.IsReadOnly)
	session.Save()

	ctx.JSON(http.StatusOK, SignInResponse{
		Code: http.StatusOK,
		User: user,
	})
}
func (h *AuthHandler) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{})
}
