package service

import (
	"log"

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

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)); err != nil {
		return nil, err
	}

	log.Println(user.UserGroup)
	return user, nil
}

func (s *AuthService) IsFullAccess(id uuid.UUID) bool {
	user, err := s.userRepo.FindByID(id, true)
	if err != nil {
		return false
	}

	return user.UserGroup.IsFullAccess
}
