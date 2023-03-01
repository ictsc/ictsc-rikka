package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type RecreateHandler struct {
	UserRepository     *repository.UserRepository
	ProblemService     *service.ProblemService
	RecreateController *controller.RecreateController
}

func NewRecreateHandler(r *gin.RouterGroup, userRepo repository.UserRepository, problemService *service.ProblemService, recreateController *controller.RecreateController) {
	handler := RecreateHandler{
		RecreateController: recreateController,
	}

	route := r.Group("/recreate")
	{
		authed := route.Group("")

		authed.Use(middleware.Auth(userRepo))

		authed.GET("/:probcode", handler.GetStatus)
		authed.POST("/:probcode", handler.CreateRecreateRequest)
	}
}

func (rh *RecreateHandler) GetStatus(ctx *gin.Context) {
	group := ctx.MustGet("group").(*entity.UserGroup)
	probcode := ctx.Param("probcode")
	b, err := rh.RecreateController.GetStatus(group, probcode)
	if err != nil {
		response.JSONRaw(ctx, http.StatusServiceUnavailable, err.Error(), nil, err)
		return
	}
	response.JSONRaw(ctx, http.StatusOK, "", b, nil)

}

func (rh *RecreateHandler) CreateRecreateRequest(ctx *gin.Context) {
	group := ctx.MustGet("group").(*entity.UserGroup)
	probcode := ctx.Param("probcode")
	b, err := rh.RecreateController.CreateRequest(group, probcode)
	if err != nil {
		response.JSONRaw(ctx, http.StatusServiceUnavailable, err.Error(), nil, err)
		return
	}
	response.JSONRaw(ctx, http.StatusOK, "", b, nil)

}
