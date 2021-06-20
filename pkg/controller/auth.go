package controller

import (
	"fmt"

	"github.com/gin-contrib/sessions"
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

func (c *AuthController) Self(session sessions.Session) (*SelfResponse, error) {
	idString, ok := session.Get("id").(string)
	if !ok || idString == "" {
		return nil, fmt.Errorf("invalid session")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	user, err := c.userService.FindByID(id)
	return &SelfResponse{
		User: user,
	}, err
}

type SignInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInResponse struct {
	User *entity.User `json:"user"`
}

func (c *AuthController) SignIn(req *SignInRequest, session sessions.Session) (*SignInResponse, error) {
	user, err := c.authService.SignIn(req.Name, req.Password, session)
	return &SignInResponse{
		User: user,
	}, err
}

func (c *AuthController) SignOut(session sessions.Session) error {
	session.Clear()
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})
	session.Save()
	return nil
}
