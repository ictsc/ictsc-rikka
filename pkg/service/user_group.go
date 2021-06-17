package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type UserGroupService struct {
	userGroupRepo repository.UserGroupRepository
}

func NewUserGroupService(userGroupRepo repository.UserGroupRepository) *UserGroupService {
	return &UserGroupService{
		userGroupRepo: userGroupRepo,
	}
}

func (s *UserGroupService) Create(userGroup *entity.UserGroup) (*entity.UserGroup, error) {
	return s.userGroupRepo.Create(userGroup)
}

func (s *UserGroupService) FindByID(id uuid.UUID) (*entity.UserGroup, error) {
	return s.userGroupRepo.FindByID(id)
}
