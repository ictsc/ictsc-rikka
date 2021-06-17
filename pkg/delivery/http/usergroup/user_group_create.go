package usergroup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type CreateRequest struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	IsFullAccess bool   `json:"is_full_access"`
}

type CreateResponse struct {
	Code           int               `json:"code"`
	InvitationCode string            `json:"invitation_code"`
	UserGroup      *entity.UserGroup `json:"user_group"`
}

func (h *UserGroupHandler) Create(ctx *gin.Context) {
	req := CreateRequest{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CreateResponse{
			Code: http.StatusBadRequest,
		})
		return
	}

	invitationCode, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, CreateResponse{
			Code: http.StatusBadRequest,
		})
		return
	}

	digest, err := bcrypt.GenerateFromPassword([]byte(invitationCode.String()), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CreateResponse{
			Code: http.StatusInternalServerError,
		})
		return
	}

	group, err := h.userGroupService.Create(&entity.UserGroup{
		Name:                 req.Name,
		Organization:         req.Organization,
		InvitationCodeDigest: string(digest),
		IsFullAccess:         req.IsFullAccess,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CreateResponse{
			Code: http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, CreateResponse{
		Code:           http.StatusCreated,
		InvitationCode: invitationCode.String(),
		UserGroup:      group,
	})
}
