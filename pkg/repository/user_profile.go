package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type UserProfileRepository interface {
	Create(profile *entity.UserProfile) (*entity.UserProfile, error)
	FindByUserID(userID uuid.UUID) (*entity.UserProfile, error)
	Update(profile *entity.UserProfile) (*entity.UserProfile, error)
}
