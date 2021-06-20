package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserGroupService struct {
	userGroupRepo repository.UserGroupRepository
}

func NewUserGroupService(userGroupRepo repository.UserGroupRepository) *UserGroupService {
	return &UserGroupService{
		userGroupRepo: userGroupRepo,
	}
}

func (s *UserGroupService) Create(name, organization, invitationCode string, isFullAccess bool) (*entity.UserGroup, error) {
	digest, err := bcrypt.GenerateFromPassword([]byte(invitationCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.userGroupRepo.Create(&entity.UserGroup{
		Name:                 name,
		Organization:         organization,
		InvitationCodeDigest: string(digest),
		IsFullAccess:         isFullAccess,
	})
}

func (s *UserGroupService) FindByID(id uuid.UUID) (*entity.UserGroup, error) {
	return s.userGroupRepo.FindByID(id)
}
