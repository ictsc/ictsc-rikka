package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type FindByIDRequest struct {
}

type FindByIDResponse struct {
	Code int          `json:"code"`
	User *entity.User `json:"user,omitempty"`
}

func (h *UserHandler) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, FindByIDResponse{
			Code: http.StatusBadRequest,
		})
		return
	}

	u, err := h.userService.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CreateResponse{
			Code: http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, CreateResponse{
		Code: http.StatusCreated,
		User: u,
	})
}
