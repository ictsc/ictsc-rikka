package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type SelfResponse struct {
	Code int          `json:"code"`
	User *entity.User `json:"user,omitempty"`
}

func (h *AuthHandler) Self(ctx *gin.Context) {

	user, ok := ctx.MustGet("user").(*entity.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, SelfResponse{
			Code: http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, SelfResponse{
		Code: http.StatusCreated,
		User: user,
	})
}
