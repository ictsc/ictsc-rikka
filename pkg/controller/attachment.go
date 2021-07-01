package controller

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
	"github.com/minio/minio-go/v7"
)

type AttachmentController struct {
	attachmentService *service.AttachmentService
}

func NewAttachmentController(attachmentService *service.AttachmentService) *AttachmentController {
	return &AttachmentController{
		attachmentService: attachmentService,
	}
}

func (c *AttachmentController) Upload(attachment *entity.Attachment) error {
	return c.attachmentService.Upload(attachment)
}
func (c *AttachmentController) Delete(id uuid.UUID) error {
	return c.attachmentService.Delete(id)
}
func (c *AttachmentController) Get(id uuid.UUID) (*minio.Object, error) {
	return c.attachmentService.Get(id)
}
func (c *AttachmentController) GetAll() ([]*minio.ObjectInfo, error) {
	return c.attachmentService.GetAll()
}
