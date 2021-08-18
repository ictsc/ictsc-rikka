package mariadb

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{
		db: db,
	}
}

func (r *UserProfileRepository) Create(profile *entity.UserProfile) (*entity.UserProfile, error) {
	if err := r.db.Create(profile).Error; err != nil {
		return nil, err
	}
	return r.FindByUserID(profile.UserID)
}

func (r *UserProfileRepository) FindByUserID(userID uuid.UUID) (*entity.UserProfile, error) {
	res := &entity.UserProfile{}
	err := r.db.Where("user_id", userID).First(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *UserProfileRepository) Update(profile *entity.UserProfile) (*entity.UserProfile, error) {
	if err := r.db.Save(profile).Error; err != nil {
		return nil, err
	}
	return r.FindByUserID(profile.UserID)
}
