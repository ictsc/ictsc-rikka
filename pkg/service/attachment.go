package service

import (
	"io"

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

func (s *AttachmentService) Create(attachment *entity.Attachment, reader io.Reader) error {
	id, err := s.attachmentRepo.Create(attachment)
	if err != nil {
		return err
	}
	err = s.s3Repo.Create(id, reader)
	if err != nil {
		return err
	}
	return nil
}
func (s *AttachmentService) Delete(id string) error {
	if err := s.s3Repo.Delete(id); err != nil {
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
