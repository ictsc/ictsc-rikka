package controller

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupController struct {
	userGroupService service.UserGroupService
}

func NewUserGroupController(userGroupService service.UserGroupService) *UserGroupController {
	return &UserGroupController{
		userGroupService: userGroupService,
	}
}

type CreateUserGroupRequest struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	IsFullAccess bool   `json:"is_full_access"`
}

type CreateUserGroupResponse struct {
	InvitationCode string            `json:"invitation_code"`
	UserGroup      *entity.UserGroup `json:"user_group"`
}

func (c *UserGroupController) Create(req *CreateUserGroupRequest) (*CreateUserGroupResponse, error) {
	invitationCode, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	userGroup, err := c.userGroupService.Create(req.Name, req.Organization, invitationCode.String(), req.IsFullAccess)
	return &CreateUserGroupResponse{
		InvitationCode: invitationCode.String(),
		UserGroup:      userGroup,
	}, err
}
