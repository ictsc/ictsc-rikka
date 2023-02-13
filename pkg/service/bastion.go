package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type BastionService struct {
	bastionRepo repository.BastionRepository
}

func NewBastionService(bastionRepo repository.BastionRepository) *BastionService {
	return &BastionService{
		bastionRepo: bastionRepo,
	}
}

func (s *BastionService) Create(userGroupId uuid.UUID, user string, password string, host string, port int) (*entity.Bastion, error) {
	return s.bastionRepo.Create(&entity.Bastion{
		UserGroupID: userGroupId,
		User:        user,
		Password:    password,
		Host:        host,
		Port:        port,
	})
}

func (s *BastionService) FindByID(id uuid.UUID) (*entity.Bastion, error) {
	return s.bastionRepo.FindByID(id)
}
