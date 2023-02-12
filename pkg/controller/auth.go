package controller

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type AuthController struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthController(authService *service.AuthService, userService *service.UserService) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

type SelfResponse struct {
	User *entity.User `json:"user"`
}

func (c *AuthController) Self(userID string) (*SelfResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	user, err := c.userService.FindMe(id)
	return &SelfResponse{
		User: user,
	}, err
}

type SignInRequest struct {
	Name     string `json:"name" binding:"required,max=32,excludesall=\r\n"`
	Password string `json:"password" binding:"required,min=8,max=40,excludesall=\r\n"`
}

type SignInResponse struct {
	User *entity.User `json:"user"`
}

func (c *AuthController) SignIn(req *SignInRequest) (*SignInResponse, error) {
	user, err := c.authService.SignIn(req.Name, req.Password)
	return &SignInResponse{
		User: user,
	}, err
}
