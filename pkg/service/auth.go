package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) SignIn(name, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByName(name, true)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("no such user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) IsFullAccess(id uuid.UUID) bool {
	user, err := s.userRepo.FindByID(id, true, false)
	if err != nil || user == nil {
		return false
	}

	return user.UserGroup.IsFullAccess
}
