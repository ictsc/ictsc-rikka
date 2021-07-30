package controller

import (
	"io"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type AttachmentController struct {
	attachmentService *service.AttachmentService
}

func NewAttachmentController(attachmentService *service.AttachmentService) *AttachmentController {
	return &AttachmentController{
		attachmentService: attachmentService,
	}
}

func (c *AttachmentController) Upload(attachment *entity.Attachment, reader io.Reader) (*entity.Attachment, error) {
	out, err := c.attachmentService.Create(attachment, reader)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *AttachmentController) Delete(id uuid.UUID) error {
	return c.attachmentService.Delete(id)
}
func (c *AttachmentController) Get(id string) (io.Reader, error) {
	return c.attachmentService.Get(id)
}
