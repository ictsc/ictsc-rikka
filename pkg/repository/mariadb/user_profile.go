package mariadb

import (
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

func (r *UserProfileRepository) FindByUserID(userID uuid.UUID) (*entity.UserProfile, error) {
	res := &entity.UserProfile{}
	err := r.db.Where("user_id", userID).First(res).Error
	return res, err
}

func (r *UserProfileRepository) UpdateOrCreate(profile *entity.UserProfile) (*entity.UserProfile, error) {
	if _, err := r.FindByUserID(profile.UserID); err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return nil, err
		}
		if err := r.db.Create(profile).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Save(profile).Error; err != nil {
			return nil, err
		}
	}
	return r.FindByUserID(profile.UserID)
}
