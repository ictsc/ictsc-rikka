package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
	"time"
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
	now := time.Now()
	notice.CreatedAt = now
	notice.UpdatedAt = now

	err := r.db.Create(notice).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(notice.ID)
}

func (r *NoticeRepository) GetAll() ([]*entity.Notice, error) {
	notice := make([]*entity.Notice, 0)
	err := r.db.Where("draft = ?", false).Order("created_at desc").Find(&notice).Error
	return notice, err
}

func (r *NoticeRepository) FindByID(id uuid.UUID) (*entity.Notice, error) {
	res := &entity.Notice{}
	err := r.db.First(res, id).Error
	return res, err
}

func (r *NoticeRepository) Update(notice *entity.Notice, skipUpdatedAt bool) (*entity.Notice, error) {
	if !skipUpdatedAt {
		notice.UpdatedAt = time.Now()
	}

	err := r.db.Save(notice).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(notice.ID)
}
