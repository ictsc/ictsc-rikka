package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
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
		authed.Use(middleware.Auth(userRepo))
		authed.GET("", handler.FindByProblemAndUserGroup)
		authed.POST("", handler.Create)
		authed.GET("/:answer_id", handler.FindByID)

		privileged := answers.Group("")
		privileged.Use(middleware.AuthIsFullAccess(userRepo))
		privileged.PATCH("/:answer_id", handler.Update)
	}
}

func (h *AnswerHandler) Create(ctx *gin.Context) {
	req := &controller.CreateAnswerRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}

	group := ctx.MustGet("group").(*entity.UserGroup)
	problem_id := ctx.Param("id")
	res, err := h.answerController.Create(group, problem_id, req)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *AnswerHandler) Update(ctx *gin.Context) {
	id := ctx.Param("answer_id")
	req := &controller.UpdateAnswerRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.Error(error.NewBadRequestError(err.Error()))
		return
	}

	res, err := h.answerController.Update(id, req)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusAccepted, "", res, nil)
}

func (h *AnswerHandler) FindByID(ctx *gin.Context) {
	group := ctx.MustGet("group").(*entity.UserGroup)
	id := ctx.Param("answer_id")

	res, err := h.answerController.FindByID(group, id)
	if err != nil {
		ctx.Error(error.NewInternalServerError(err))
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AnswerHandler) FindByProblemAndUserGroup(ctx *gin.Context) {
	group := ctx.MustGet("group").(*entity.UserGroup)
	probid := ctx.Param("id")

	if group.IsFullAccess {
		res, err := h.answerController.FindByProblem(group, probid, nil)
		if err != nil {
			ctx.Error(error.NewInternalServerError(err))
			return
		}
		response.JSON(ctx, http.StatusOK, "", res, nil)
	} else {
		res, err := h.answerController.FindByProblemAndUserGroup(group, probid)
		if err != nil {
			ctx.Error(error.NewInternalServerError(err))
			return
		}
		response.JSON(ctx, http.StatusOK, "", res, nil)
	}

}
