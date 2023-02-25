package mariadb

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type BastionRepository struct {
	db *gorm.DB
}

func NewBastionRepository(db *gorm.DB) *BastionRepository {
	return &BastionRepository{
		db: db,
	}
}

func (r *BastionRepository) Create(bastion *entity.Bastion) (*entity.Bastion, error) {
	if err := r.db.Create(bastion).Error; err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *BastionRepository) FindByID(id uuid.UUID) (*entity.Bastion, error) {
	res := &entity.Bastion{}
	err := r.db.First(res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}
