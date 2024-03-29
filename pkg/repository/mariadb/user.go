package mariadb

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *entity.User) (*entity.User, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return r.FindByID(user.ID, false, false)
}

func (r *UserRepository) FindByID(id uuid.UUID, isPreload bool, isPreloadBastionData bool) (*entity.User, error) {
	db := r.db
	res := &entity.User{}
	if isPreload {
		if isPreloadBastionData {
			db = db.Preload("UserGroup.Bastion").Preload("UserProfile")
		} else {
			db = preload(db)
		}
	}
	err := db.First(res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *UserRepository) FindByUserGroupID(id uuid.UUID) ([]*entity.User, error) {
	res := make([]*entity.User, 0)
	err := preload(r.db).Where("user_group_id", id).Find(&res).Error
	return res, err
}

func (r *UserRepository) FindByName(name string, isPreload bool) (*entity.User, error) {
	db := r.db
	res := &entity.User{}
	if isPreload {
		db = preload(db)
	}
	err := db.Where("name", name).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *UserRepository) Update(user *entity.User) (*entity.User, error) {
	user.UpdatedAt = time.Now()

	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return r.FindByID(user.ID, true, false)
}

func preload(db *gorm.DB) *gorm.DB {
	return db.Preload("UserGroup").Preload("UserProfile")
}
