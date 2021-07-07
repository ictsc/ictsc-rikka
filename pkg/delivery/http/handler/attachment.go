package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
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
	user := ctx.MustGet("user").(*entity.User)
	file, err := ctx.FormFile("file")
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
	}
	reader, err := file.Open()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
	}
	attachment := &entity.Attachment{
		User: user.ID,
	}
	err = h.attachmentController.Upload(attachment, reader)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
	}

	response.JSON(ctx, http.StatusCreated, "", "", nil)
}

func (h *AttachmentHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := h.attachmentController.Get(id)
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

	err := h.attachmentController.Delete(id)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}
	response.JSON(ctx, http.StatusNoContent, "", nil, nil)
}
