package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type SelfResponse struct {
	User *entity.User `json:"user,omitempty"`
}

func (h *AuthHandler) Self(ctx *gin.Context) {

	user, ok := ctx.MustGet("user").(*entity.User)
	if !ok {
		response.JSON(ctx, http.StatusInternalServerError, "", nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", SelfResponse{
		User: user,
	}, nil)
}
