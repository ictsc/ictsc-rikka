package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

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
	session.Save()

	ctx.JSON(http.StatusOK, SignInResponse{
		Code: http.StatusOK,
		User: user,
	})
}
