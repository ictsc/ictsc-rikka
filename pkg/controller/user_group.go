package controller

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupController struct {
	userGroupService *service.UserGroupService
}

func NewUserGroupController(userGroupService *service.UserGroupService) *UserGroupController {
	return &UserGroupController{
		userGroupService: userGroupService,
	}
}

type CreateUserGroupRequest struct {
	Name           string `json:"name"`
	Organization   string `json:"organization"`
	InvitationCode string `json:"invitation_code"`
	IsFullAccess   bool   `json:"is_full_access"`
}

type CreateUserGroupResponse struct {
	UserGroup *entity.UserGroup `json:"user_group"`
}

func (c *UserGroupController) Create(req *CreateUserGroupRequest) (*CreateUserGroupResponse, error) {
	userGroup, err := c.userGroupService.Create(req.Name, req.Organization, req.InvitationCode, req.IsFullAccess)
	return &CreateUserGroupResponse{
		UserGroup: userGroup,
	}, err
}
