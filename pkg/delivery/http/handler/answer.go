package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
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
		response.JSON(ctx, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	groupuuid := ctx.MustGet("group").(*entity.UserGroup).ID
	problem_id := ctx.Param("id")
	res, err := h.answerController.Create(problem_id, groupuuid, req)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusCreated, "", res, nil)
}

func (h *AnswerHandler) Update(ctx *gin.Context) {
	id := ctx.Param("answer_id")
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

func (h *AnswerHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("answer_id")

	res, err := h.answerController.FindByID(id)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AnswerHandler) FindByProblemAndUserGroup(ctx *gin.Context) {
	probid := ctx.Param("id")

	is_full_access := ctx.MustGet("is_full_access").(bool)

	if is_full_access {
		userGroupID := ctx.Param("user_group_id")
		res, err := h.answerController.FindByProblem(probid, userGroupID)
		if err != nil {
			response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
			return
		}
		response.JSON(ctx, http.StatusOK, "", res, nil)
	} else {
		userGroupID := ctx.MustGet("group").(*entity.UserGroup).ID
		res, err := h.answerController.FindByProblemAndUserGroup(probid, userGroupID)
		if err != nil {
			response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
			return
		}
		response.JSON(ctx, http.StatusOK, "", res, nil)
	}

}

func (h *AnswerHandler) GetAll(ctx *gin.Context) {
	res, err := h.answerController.GetAll()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusOK, "", res, nil)
}

func (h *AnswerHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("answer_id")

	err := h.answerController.Delete(id)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response.JSON(ctx, http.StatusNoContent, "", nil, nil)
}
