package controller

import (
	"io"

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

func (c *AttachmentController) Upload(attachment *entity.Attachment, reader io.Reader) error {
	return c.attachmentService.Create(attachment, reader)
}
func (c *AttachmentController) Delete(id string) error {
	return c.attachmentService.Delete(id)
}
func (c *AttachmentController) Get(id string) (io.Reader, error) {
	return c.attachmentService.Get(id)
}
func (c *AttachmentController) GetAll() ([]*minio.ObjectInfo, error) {
	return c.attachmentService.GetAll()
}
