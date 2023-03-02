package controller

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserGroupController struct {
	userService      *service.UserService
	userGroupService *service.UserGroupService
	bastionService   *service.BastionService
}

func NewUserGroupController(userService *service.UserService, userGroupService *service.UserGroupService, bastionService *service.BastionService) *UserGroupController {
	return &UserGroupController{
		userService:      userService,
		userGroupService: userGroupService,
		bastionService:   bastionService,
	}
}

type ListParticipatesMember struct {
	DisplayName string              `json:"display_name"`
	Profile     *entity.UserProfile `json:"profile"`
}

type ListParticipates struct {
	Name         string                    `json:"name"`
	Organization string                    `json:"organization"`
	Members      []*ListParticipatesMember `json:"members"`
}

func (c *UserGroupController) ListParticipates() ([]*ListParticipates, error) {
	userGroups, err := c.userGroupService.List()
	if err != nil {
		return nil, err
	}

	resp := make([]*ListParticipates, 0, len(userGroups))
	for _, userGroup := range userGroups {
		users, err := c.userService.FindByUserGroupID(userGroup.ID)
		if err != nil {
			return nil, err
		}

		members := make([]*ListParticipatesMember, 0, len(users))
		for _, user := range users {
			members = append(members, &ListParticipatesMember{
				DisplayName: user.DisplayName,
				Profile:     user.UserProfile,
			})
		}

		resp = append(resp, &ListParticipates{
			Name:         userGroup.Name,
			Organization: userGroup.Organization,
			Members:      members,
		})
	}

	return resp, err
}

type CreateUserGroupRequest struct {
	Name            string `json:"name"`
	Organization    string `json:"organization"`
	InvitationCode  string `json:"invitation_code"`
	IsFullAccess    bool   `json:"is_full_access"`
	BastionUser     string `json:"bastion_user"`
	BastionPassword string `json:"bastion_password"`
	BastionHost     string `json:"bastion_host"`
	BastionPort     int    `json:"bastion_port"`
	TeamID          string `json:"team_id"`
}

type CreateUserGroupResponse struct {
	UserGroup *entity.UserGroup `json:"user_group"`
}

func (c *UserGroupController) Create(req *CreateUserGroupRequest) (*CreateUserGroupResponse, error) {
	userGroup, err := c.userGroupService.Create(req.Name, req.Organization, req.InvitationCode, req.IsFullAccess, req.TeamID)
	if err != nil {
		return nil, err
	}

	_, err = c.bastionService.Create(userGroup.ID, req.BastionUser, req.BastionPassword, req.BastionHost, req.BastionPort)
	if err != nil {
		return nil, err
	}

	return &CreateUserGroupResponse{
		UserGroup: userGroup,
	}, err
}
