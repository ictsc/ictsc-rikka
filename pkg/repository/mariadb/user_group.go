package mariadb

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type UserGroupRepository struct {
	db *gorm.DB
}

func NewUserGroupRepository(db *gorm.DB) *UserGroupRepository {
	return &UserGroupRepository{
		db: db,
	}
}

func (r *UserGroupRepository) Create(userGroup *entity.UserGroup) (*entity.UserGroup, error) {
	now := time.Now()
	userGroup.CreatedAt = now
	userGroup.UpdatedAt = now

	err := r.db.Create(userGroup).Error
	if err != nil {
		return nil, err
	}

	return r.FindByID(userGroup.ID)
}

func (r *UserGroupRepository) GetAll() ([]*entity.UserGroup, error) {
	userGroups := make([]*entity.UserGroup, 0)
	err := r.db.Find(&userGroups).Error
	return userGroups, err
}

func (r *UserGroupRepository) FindByID(id uuid.UUID) (*entity.UserGroup, error) {
	res := &entity.UserGroup{}
	err := r.db.First(res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *UserGroupRepository) FindByName(name string) (*entity.UserGroup, error) {
	res := &entity.UserGroup{}
	err := r.db.Where("name", name).First(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}
