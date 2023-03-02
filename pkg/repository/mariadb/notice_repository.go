package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type NoticeRepository struct {
	db *gorm.DB
}

func NewNoticeRepository(db *gorm.DB) *NoticeRepository {
	return &NoticeRepository{
		db: db,
	}
}

func (r *NoticeRepository) Create(notice *entity.Notice) (*entity.Notice, error) {
	err := r.db.Create(notice).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(notice.ID)
}

func (r *NoticeRepository) GetAll() ([]*entity.Notice, error) {
	notice := make([]*entity.Notice, 0)
	err := r.db.Find(&notice).Error
	return notice, err
}

func (r *NoticeRepository) FindByID(id uuid.UUID) (*entity.Notice, error) {
	res := &entity.Notice{}
	err := r.db.First(res, id).Error
	return res, err
}

func (r *NoticeRepository) Update(notice *entity.Notice) (*entity.Notice, error) {
	err := r.db.Save(notice).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(notice.ID)
}
