package mariadb

import (
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
	err := r.db.Create(userGroup).Error
	if err != nil {
		return nil, err
	}

	res := &entity.UserGroup{}
	err = r.db.First(res, userGroup.ID).Error
	return res, err
}

func (r *UserGroupRepository) FindByID(id uuid.UUID) (*entity.UserGroup, error) {
	res := &entity.UserGroup{}
	err := r.db.First(res, id).Error
	return res, err
}

func (r *UserGroupRepository) FindByName(name string) (*entity.UserGroup, error) {
	res := &entity.UserGroup{}
	err := r.db.Where("name", name).First(res).Error
	return res, err
}
