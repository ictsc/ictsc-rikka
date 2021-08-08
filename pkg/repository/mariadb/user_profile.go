package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type UserProfileRepository struct {
	*db
}

func NewUserProfileRepository(config *MariaDBConfig) *UserProfileRepository {
	return &UserProfileRepository{
		db: newDB(config),
	}
}

func (r *UserProfileRepository) FindByUserID(userID uuid.UUID) (*entity.UserProfile, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.UserProfile{}
	err = db.Where("user_id", userID).First(res).Error
	return res, err
}

func (r *UserProfileRepository) UpdateOrCreate(profile *entity.UserProfile) (*entity.UserProfile, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err := r.FindByUserID(profile.UserID); err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return nil, err
		}
		if err := db.Create(profile).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Save(profile).Error; err != nil {
			return nil, err
		}
	}
	return r.FindByUserID(profile.UserID)
}
