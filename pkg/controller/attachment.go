package controller

import (
	"io"

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

type UploadAttachmentRequest struct {
	ID     uuid.UUID
	Reader io.Reader
}
type UploadAttachmentResponse struct {
	Attachment *entity.Attachment
}

func (c *AttachmentController) Upload(req *UploadAttachmentRequest) (*UploadAttachmentResponse, error) {
	//id, _ := uuid.NewRandom()
	attachment, err := c.attachmentService.Upload(&service.UploadAttachmentRequest{
		Reader: req.Reader,
	})
	if err != nil {
		return nil, err
	}

	return &UploadAttachmentResponse{
		Attachment: attachment,
	}, nil
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
