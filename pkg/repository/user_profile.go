package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type UserProfileRepository interface {
	FindByUserID(userID uuid.UUID) (*entity.UserProfile, error)
	UpdateOrCreate(profile *entity.UserProfile) (*entity.UserProfile, error)
}
