package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
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
	User *entity.User `json:"user,omitempty"`
}

func (h *UserHandler) Create(ctx *gin.Context) {
	req := CreateRequest{}
	if err := ctx.Bind(&req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	digest, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	userGroupID, err := uuid.Parse(req.UserGroupID)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	u, err := h.userService.Create(&entity.User{
		Name:           req.Name,
		DisplayName:    req.Name,
		PasswordDigest: string(digest),
		UserGroupID:    userGroupID,
	}, req.InvitationCode)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", CreateResponse{
		User: u,
	}, nil)
}
