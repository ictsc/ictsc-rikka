package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/minio/minio-go/v7"
)

type AttachmentService struct {
	attachmentRepo repository.AttachmentRepository
	s3Repo         repository.S3Repository
}

type DeleteAttachmentRequest struct {
	Name string
}
type GetAttachmentRequest struct {
	Name string
}

func NewAttachmentService(attachmentRepo repository.AttachmentRepository, s3Repo repository.S3Repository) *AttachmentService {
	return &AttachmentService{
		attachmentRepo: attachmentRepo,
		s3Repo:         s3Repo,
	}
}

func (s *AttachmentService) Upload(attachment *entity.Attachment) error {
	id, _ := uuid.NewRandom()
	attachment.ID = id
	if err := s.s3Repo.Create(attachment); err != nil {
		return err
	}
	if err := s.attachmentRepo.Upload(attachment); err != nil {
		return err
	}
	return nil
}
func (s *AttachmentService) Delete(id uuid.UUID) error {
	if err := s.s3Repo.Delete(id); err != nil {
		return err
	}
	if err := s.attachmentRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
func (s *AttachmentService) Get(id uuid.UUID) (*minio.Object, error) {
	obj, err := s.attachmentRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (s *AttachmentService) GetAll() ([]*minio.ObjectInfo, error) {
	return s.attachmentRepo.GetAll()
}
