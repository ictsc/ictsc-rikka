package mariadb

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type DBAttachmentRepository struct {
	db *gorm.DB
}

func NewDBAttachmentRepository(db *gorm.DB) *DBAttachmentRepository {
	return &DBAttachmentRepository{
		db: db,
	}
}

func (r *DBAttachmentRepository) Create(attachment *entity.Attachment) (*entity.Attachment, error) {
	err := r.db.Create(attachment).Error
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
func (r *DBAttachmentRepository) Delete(Attachment *entity.Attachment) error {
	return r.db.Delete(Attachment, Attachment.ID).Error
}
func (r *DBAttachmentRepository) Get(req *entity.Attachment) (*entity.Attachment, error) {
	attachment := &entity.Attachment{}
	err := r.db.Where("id", req.ID).Find(attachment).Error
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
func (r *DBAttachmentRepository) GetAll() ([]*entity.Attachment, error) {
	attachments := make([]*entity.Attachment, 0)
	err := r.db.Find(attachments).Error
	return attachments, err
}
