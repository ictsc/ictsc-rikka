package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type NoticeRepository interface {
	Create(notice *entity.Notice) (*entity.Notice, error)
	GetAll() ([]*entity.Notice, error)
	FindByID(id uuid.UUID) (*entity.Notice, error)
	Update(notice *entity.Notice) (*entity.Notice, error)
}
