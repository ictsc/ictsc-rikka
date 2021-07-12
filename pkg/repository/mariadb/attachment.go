package mariadb

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{
		db: db,
	}
}

func (r *AttachmentRepository) Create(attachment *entity.Attachment) (string, error) {
	err := r.db.Create(attachment).Error
	return attachment.Base.ID.String(), err
}
func (r *AttachmentRepository) Delete(id string) error {
	attachment := &entity.Attachment{}
	return r.db.Delete(attachment, id).Error
}
func (r *AttachmentRepository) Get(id string) (*entity.Attachment, error) {
	attachment := &entity.Attachment{}
	err := r.db.Where("id", id).Find(attachment).Error
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
func (r *AttachmentRepository) GetAll() ([]*entity.Attachment, error) {
	attachments := make([]*entity.Attachment, 0)
	err := r.db.Find(attachments).Error
	return attachments, err
}
