package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"net/http"
)

type AnswerHandler struct {
	answerController *controller.AnswerController
}

func NewAnswerHandler(r *gin.RouterGroup, userRepo repository.UserRepository, answerController *controller.AnswerController) {
	handler := AnswerHandler{
		answerController: answerController,
	}

	answers := r.Group("/answers")
	{
		authed := answers.Group("")
		authed.Use(middleware.AuthIsFullAccess(userRepo))
		privileged := answers.Group("")
		privileged.Use(middleware.AuthIsFullAccess(userRepo))

		authed.GET("", handler.GetAll)
		authed.GET("/:id", handler.Find)

		privileged.POST("", handler.Create)
		privileged.PUT("/:id", handler.Update)
		privileged.DELETE("/:id", handler.Delete)
	}
}

func (h *AnswerHandler) Create(ctx *gin.Context) {
	req := &controller.CreateAnswerRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	res, err := h.answerController.Create(req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *AnswerHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	req := &controller.UpdateAnswerRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	res, err := h.answerController.Update(id, req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusAccepted, "", res, nil)
}

func (h *AnswerHandler) Find(ctx *gin.Context) {
	id := ctx.Param("id")
	metadataOnly := ctx.Query("metadata_only") != ""

	res, err := h.answerController.FindByID(id, metadataOnly)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AnswerHandler) GetAll(ctx *gin.Context) {
	metadataOnly := ctx.Query("metadata_only") != ""

	res, err := h.answerController.GetAll(metadataOnly)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AnswerHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.answerController.Delete(id)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusNoContent, "", nil, nil)
}
