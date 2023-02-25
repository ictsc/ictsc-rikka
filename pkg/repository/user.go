package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByID(id uuid.UUID, isPreload bool, isPreloadBastionData bool) (*entity.User, error)
	FindByUserGroupID(id uuid.UUID) ([]*entity.User, error)
	FindByName(name string, isPreload bool) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
}
