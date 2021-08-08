package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	*db
}

func NewUserRepository(config *MariaDBConfig) *UserRepository {
	return &UserRepository{
		db: newDB(config),
	}
}

func (r *UserRepository) Create(user *entity.User) (*entity.User, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	res := &entity.User{}
	err = db.First(res, user.ID).Error
	return res, err
}

func (r *UserRepository) FindByID(id uuid.UUID, isPreload bool) (*entity.User, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.User{}
	if isPreload {
		db = preload(db)
	}
	err = db.First(res, id).Error
	return res, err
}

func (r *UserRepository) FindByName(name string, isPreload bool) (*entity.User, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.User{}
	if isPreload {
		db = preload(db)
	}
	err = db.Where("name", name).First(&res).Error
	return res, err
}

func (r *UserRepository) Update(user *entity.User) (*entity.User, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := db.Save(user).Error; err != nil {
		return nil, err
	}
	return r.FindByID(user.ID, true)
}

func preload(db *gorm.DB) *gorm.DB {
	return db.Preload("UserGroup").Preload("UserProfile")
}
