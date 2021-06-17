package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type CreateRequest struct {
	Name           string `json:"name"`
	Password       string `json:"password"`
	UserGroupID    string `json:"user_group_id"`
	InvitationCode string `json:"invitation_code"`
}

type CreateResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	User    *entity.User `json:"user,omitempty"`
}

func (h *UserHandler) Create(ctx *gin.Context) {
	req := CreateRequest{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CreateResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	digest, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CreateResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	userGroupID, err := uuid.Parse(req.UserGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CreateResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	u, err := h.userService.Create(&entity.User{
		Name:           req.Name,
		DisplayName:    req.Name,
		PasswordDigest: string(digest),
		UserGroupID:    userGroupID,
	}, req.InvitationCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CreateResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, CreateResponse{
		Code: http.StatusCreated,
		User: u,
	})
}
