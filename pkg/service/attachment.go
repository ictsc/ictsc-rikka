package service

import (
	"io"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/minio/minio-go/v7"
)

type AttachmentService struct {
	attachmentRepo repository.AttachmentRepository
}
type UploadAttachmentRequest struct {
	Name   string
	Reader io.Reader
}
type DeleteAttachmentRequest struct {
	Name string
}
type GetAttachmentRequest struct {
	Name string
}

func NewAttachmentService(attachmentRepo repository.AttachmentRepository) *AttachmentService {
	return &AttachmentService{
		attachmentRepo: attachmentRepo,
	}
}

func (s *AttachmentService) Upload(req *UploadAttachmentRequest) (*entity.Attachment, error) {
	attachment := &entity.Attachment{
		Name:   req.Name,
		Reader: req.Reader,
	}
	return s.attachmentRepo.Upload(attachment)

}
func (s *AttachmentService) Delete(id uuid.UUID) error {

	return s.attachmentRepo.Delete(id)
}
func (s *AttachmentService) Get(id uuid.UUID) (*minio.Object, error) {
	return s.attachmentRepo.Get(id)
}
func (s *AttachmentService) GetAll() ([]*minio.ObjectInfo, error) {
	return s.attachmentRepo.GetAll()
}
