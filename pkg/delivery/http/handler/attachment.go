package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
)

type AttachmentHandler struct {
	attachmentController *controller.AttachmentController
}

func NewAttachmentHandler(r *gin.RouterGroup, attachmentController *controller.AttachmentController) {
	handler := AttachmentHandler{
		attachmentController: attachmentController,
	}
	attachments := r.Group("/attachments")
	{
		attachments.POST("/:id", handler.Upload)
		attachments.GET("/", handler.GetAll)
		attachments.GET("/:id", handler.Get)
		attachments.DELETE("/:id", handler.Delete)
	}
}

func (h *AttachmentHandler) Upload(ctx *gin.Context) {
	reader := ctx.Request.Body
	id, _ := uuid.NewRandom()
	req := &controller.UploadAttachmentRequest{
		ID:     id,
		Reader: reader,
	}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}
	res, err := h.attachmentController.Upload(req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *AttachmentHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}
	res, err := h.attachmentController.Get(uid)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}
	response.JSON(ctx, http.StatusOK, "", res, nil)
}
func (h *AttachmentHandler) GetAll(ctx *gin.Context) {
	res, err := h.attachmentController.GetAll()
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}
	response.JSON(ctx, http.StatusOK, "", res, nil)
}
func (h *AttachmentHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	err = h.attachmentController.Delete(uid)
	response.JSON(ctx, http.StatusNoContent, "", nil, nil)
}
