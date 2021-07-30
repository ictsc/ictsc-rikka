package repository

import (
	"io"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type S3Repository interface {
	Create(id string, reader io.Reader) error
	Delete(id string) error
	Get(id string) (io.Reader, error)
}
type AttachmentRepository interface {
	Create(attachment *entity.Attachment) (*entity.Attachment, error)
	Delete(id uuid.UUID) error
	Get(id string) (*entity.Attachment, error)
}
