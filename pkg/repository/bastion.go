package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type BastionRepository interface {
	Create(bastion *entity.Bastion) (*entity.Bastion, error)
	FindByID(id uuid.UUID) (*entity.Bastion, error)
}
