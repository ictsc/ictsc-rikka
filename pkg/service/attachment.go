package service

import (
	"io"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
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

func (s *AttachmentService) Create(attachment *entity.Attachment, reader io.Reader) (*entity.Attachment, error) {
	created, err := s.attachmentRepo.Create(attachment)
	if err != nil {
		return nil, err
	}
	err = s.s3Repo.Create(created.ID.String(), reader)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *AttachmentService) Delete(id uuid.UUID) error {
	if err := s.s3Repo.Delete(id.String()); err != nil {
		return err
	}
	if err := s.attachmentRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *AttachmentService) Get(id string) (io.Reader, error) {
	obj, err := s.s3Repo.Get(id)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
