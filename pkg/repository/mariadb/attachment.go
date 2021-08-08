package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type AttachmentRepository struct {
	*db
}

func NewAttachmentRepository(config *MariaDBConfig) *AttachmentRepository {
	return &AttachmentRepository{
		db: newDB(config),
	}
}

func (r *AttachmentRepository) Create(attachment *entity.Attachment) (*entity.Attachment, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := db.Create(attachment).Error; err != nil {
		return nil, err
	}
	return r.Get(attachment.ID.String())
}

func (r *AttachmentRepository) Delete(id uuid.UUID) error {
	db, conn, err := r.init()
	if err != nil {
		return err
	}
	defer conn.Close()

	attachment := &entity.Attachment{
		Base: entity.Base{
			ID: id,
		},
	}
	return db.Delete(attachment).Error
}
func (r *AttachmentRepository) Get(id string) (*entity.Attachment, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	attachment := &entity.Attachment{}
	err = db.Where("id", id).Find(attachment).Error
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
func (r *AttachmentRepository) GetAll() ([]*entity.Attachment, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	attachments := make([]*entity.Attachment, 0)
	err = db.Find(attachments).Error
	return attachments, err
}
