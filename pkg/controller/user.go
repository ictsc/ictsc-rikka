package controller

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type CreateUserRequest struct {
	Name           string `json:"name"`
	Password       string `json:"password"`
	UserGroupID    string `json:"user_group_id"`
	InvitationCode string `json:"invitation_code"`
}

type CreateUserResponse struct {
	User *entity.User `json:"user"`
}

func (c *UserController) Create(req *CreateUserRequest) (*CreateUserResponse, error) {
	userGroupID, err := uuid.Parse(req.UserGroupID)
	if err != nil {
		return nil, err
	}

	user, err := c.userService.Create(req.Name, req.Password, userGroupID, req.InvitationCode)
	return &CreateUserResponse{
		User: user,
	}, err
}

type FindUserByIDResponse struct {
	User *entity.User `json:"user"`
}

func (c *UserController) FindByID(userID string) (*FindUserByIDResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	user, err := c.userService.FindByID(id)
	return &FindUserByIDResponse{
		User: user,
	}, err
}

type UpdateUserRequest struct {
	DisplayName      string `json:"display_name"`
	TwitterID        string `json:"twitter_id"`
	GithubID         string `json:"github_id"`
	FacebookID       string `json:"facebook_id"`
	SelfIntroduction string `json:"self_introduction"`
}

type UpdateUserResponse struct {
	User *entity.User `json:"user"`
}

func (c *UserController) Update(userID string, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	user, err := c.userService.Update(id, req.DisplayName, req.TwitterID, req.GithubID, req.FacebookID, req.SelfIntroduction)
	if err != nil {
		return nil, err
	}
	return &UpdateUserResponse{
		User: user,
	}, nil
}
