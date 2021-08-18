package mariadb

import (
	"errors"

	"github.com/google/uuid"
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

func (r *AttachmentRepository) Create(attachment *entity.Attachment) (*entity.Attachment, error) {
	if err := r.db.Create(attachment).Error; err != nil {
		return nil, err
	}
	return r.Get(attachment.ID.String())
}

func (r *AttachmentRepository) Delete(id uuid.UUID) error {
	attachment := &entity.Attachment{
		Base: entity.Base{
			ID: id,
		},
	}
	return r.db.Delete(attachment).Error
}
func (r *AttachmentRepository) Get(id string) (*entity.Attachment, error) {
	attachment := &entity.Attachment{}
	err := r.db.Where("id", id).Find(attachment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return attachment, err
}
func (r *AttachmentRepository) GetAll() ([]*entity.Attachment, error) {
	attachments := make([]*entity.Attachment, 0)
	err := r.db.Find(attachments).Error
	return attachments, err
}
