package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"net/http"
)

type ProblemHandler struct {
	problemController *controller.ProblemController
}

func NewProblemHandler(r *gin.RouterGroup, userRepo repository.UserRepository, problemController *controller.ProblemController) {
	handler := ProblemHandler{
		problemController: problemController,
	}

	problems := r.Group("/problems")
	{
		authed := problems.Group("")
		authed.Use(middleware.AuthIsFullAccess(userRepo))
		privileged := problems.Group("")
		privileged.Use(middleware.AuthIsFullAccess(userRepo))

		authed.GET("", handler.GetAll)
		authed.GET("/:id", handler.Find)

		privileged.POST("", handler.Create)
		privileged.PUT("/:id", handler.Update)
		privileged.DELETE("/:id", handler.Delete)
	}
}

func (h *ProblemHandler) Create(ctx *gin.Context) {
	req := &controller.CreateProblemRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	res, err := h.problemController.Create(req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *ProblemHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	req := &controller.UpdateProblemRequest{}
	if err := ctx.Bind(req); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	res, err := h.problemController.Update(id, req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusAccepted, "", res, nil)
}

func (h *ProblemHandler) Find(ctx *gin.Context) {
	id := ctx.Param("id")
	metadataOnly := ctx.Query("metadata_only") != ""

	res, err := h.problemController.FindByID(id, metadataOnly)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *ProblemHandler) GetAll(ctx *gin.Context) {
	metadataOnly := ctx.Query("metadata_only") != ""

	res, err := h.problemController.GetAll(metadataOnly)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *ProblemHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.problemController.Delete(id)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusNoContent, "", nil, nil)
}
