package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type AttachmentHandler struct {
	attachmentController *controller.AttachmentController
}

func NewAttachmentHandler(r *gin.RouterGroup, attachmentController *controller.AttachmentController, userRepo repository.UserRepository) {
	handler := AttachmentHandler{
		attachmentController: attachmentController,
	}
	attachments := r.Group("/attachments")
	{
		authed := attachments.Group("")
		authed.Use(middleware.Auth(userRepo))
		{
			authed.POST("", handler.Upload)
			authed.GET("/:id", handler.Get)
			authed.DELETE("/:id", handler.Delete)

		}
	}
}

func (h *AttachmentHandler) Upload(ctx *gin.Context) {

	user := ctx.MustGet("user").(*entity.User)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(error.NewBadRequestError("FormFile isn't specified"))
		return
	}
	reader, err := file.Open()
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}
	attachment := &entity.Attachment{
		UserID: user.ID,
	}
	res, err := h.attachmentController.Upload(attachment, reader)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *AttachmentHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := h.attachmentController.Get(id)
	if err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}
	io.Copy(ctx.Writer, res)
}

func (h *AttachmentHandler) Delete(ctx *gin.Context) {
	idString := ctx.Param("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}

	if err := h.attachmentController.Delete(id); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}
	response.JSON(ctx, http.StatusNoContent, "", nil, nil)
}
