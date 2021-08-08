package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type UserGroupRepository struct {
	*db
}

func NewUserGroupRepository(config *MariaDBConfig) *UserGroupRepository {
	return &UserGroupRepository{
		db: newDB(config),
	}
}

func (r *UserGroupRepository) Create(userGroup *entity.UserGroup) (*entity.UserGroup, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = db.Create(userGroup).Error
	if err != nil {
		return nil, err
	}

	res := &entity.UserGroup{}
	err = db.First(res, userGroup.ID).Error
	return res, err
}

func (r *UserGroupRepository) FindByID(id uuid.UUID) (*entity.UserGroup, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.UserGroup{}
	err = db.First(res, id).Error
	return res, err
}

func (r *UserGroupRepository) FindByName(name string) (*entity.UserGroup, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.UserGroup{}
	err = db.Where("name", name).First(res).Error
	return res, err
}
