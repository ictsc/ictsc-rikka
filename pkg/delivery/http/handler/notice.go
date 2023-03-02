package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"net/http"
)

type NoticeHandler struct {
	noticeController *controller.NoticeController
}

func NewNoticeHandler(r *gin.RouterGroup, userRepo repository.UserRepository, noticeController *controller.NoticeController) {
	handler := NoticeHandler{
		noticeController: noticeController,
	}

	notices := r.Group("/notices")
	{
		authed := notices.Group("")
		authed.Use(middleware.Auth(userRepo))

		authed.GET("", handler.GetAll)
	}
}

func (h *NoticeHandler) GetAll(ctx *gin.Context) {
	res, err := h.noticeController.GetAll()
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}
