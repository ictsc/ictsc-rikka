package repository

import (
	"io"

	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/minio/minio-go/v7"
)

type S3Repository interface {
	Create(id string, reader io.Reader) error
	Delete(id string) error
	Get(id string) (io.Reader, error)
	GetAll() ([]*minio.ObjectInfo, error)
}
type AttachmentRepository interface {
	Create(attachment *entity.Attachment) (string, error)
	Delete(id string) error
	Get(id string) (io.Reader, error)
	GetAll() ([]*minio.ObjectInfo, error)
}
