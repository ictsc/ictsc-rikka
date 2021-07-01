package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/minio/minio-go/v7"
)

type AttachmentRepository interface {
	Upload(attachment *entity.Attachment) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*minio.Object, error)
	GetAll() ([]*minio.ObjectInfo, error)
}
type S3Repository interface {
	Create(attachment *entity.Attachment) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*minio.Object, error)
	GetAll() ([]*minio.ObjectInfo, error)
}
