package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type UserGroupRepository interface {
	Create(userGroup *entity.UserGroup) (*entity.UserGroup, error)
	FindByID(id uuid.UUID) (*entity.UserGroup, error)
}
