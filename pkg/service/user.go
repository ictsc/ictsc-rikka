package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo      repository.UserRepository
	userGroupRepo repository.UserGroupRepository
}

func NewUserService(userRepo repository.UserRepository, userGroupRepo repository.UserGroupRepository) *UserService {
	return &UserService{
		userRepo:      userRepo,
		userGroupRepo: userGroupRepo,
	}
}

func (s *UserService) Create(user *entity.User, invitationCode string) (*entity.User, error) {
	group, err := s.userGroupRepo.FindByID(user.UserGroupID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(group.InvitationCodeDigest), []byte(invitationCode)); err != nil {
		return nil, err
	}

	return s.userRepo.Create(user)
}

func (s *UserService) FindByID(id uuid.UUID) (*entity.User, error) {
	return s.userRepo.FindByID(id, false)
}
