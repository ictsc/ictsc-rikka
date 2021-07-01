package mariadb

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type S3Repository struct {
	db *gorm.DB
}

func NewS3Repository(db *gorm.DB) *S3Repository {
	return &S3Repository{
		db: db,
	}
}

func (r *S3Repository) Create(attachment *entity.Attachment) error {
	return r.db.Create(attachment).Error
}
func (r *S3Repository) Delete(Attachment *entity.Attachment) error {
	return r.db.Delete(Attachment, Attachment.ID).Error
}
func (r *S3Repository) Get(req *entity.Attachment) (*entity.Attachment, error) {
	attachment := &entity.Attachment{}
	err := r.db.Where(req, req.ID).Find(attachment).Error
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
func (r *S3Repository) GetAll() ([]*entity.Attachment, error) {
	attachments := make([]*entity.Attachment, 0)
	err := r.db.Find(attachments).Error
	return attachments, err
}
