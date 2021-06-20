package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type SignInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInResponse struct {
	User *entity.User `json:"user,omitempty"`
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	req := SignInRequest{}
	if err := ctx.Bind(&req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return

	}
	user, err := h.authService.SignIn(req.Name, req.Password)
	if err != nil {
		response.JSON(ctx, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	session := sessions.Default(ctx)
	session.Set("id", user.ID.String())
	session.Save()

	response.JSON(ctx, http.StatusOK, "", SignInResponse{
		User: user,
	}, nil)
}
