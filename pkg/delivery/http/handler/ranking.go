package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type RankingHandler struct {
	rankingController *controller.RankingController
}

func NewRankingHandler(r *gin.RouterGroup, userRepo repository.UserRepository, rankingService *service.RankingService) {
	handler := RankingHandler{
		rankingController: controller.NewRankingController(rankingService),
	}

	route := r.Group("/ranking")
	route.Use(middleware.Auth(userRepo))
	{
		route.GET("/top", handler.GetTopRanking)
		route.GET("/near-me", handler.GetNearMeRanking)
	}
}

func (h *RankingHandler) GetTopRanking(ctx *gin.Context) {
	ranking, err := h.rankingController.GetTopRanking()
	if err != nil {
		ctx.Error(err)
		return
	}

	response.JSON(ctx, http.StatusOK, "", ranking, nil)
}

func (h *RankingHandler) GetNearMeRanking(ctx *gin.Context) {
	user, ok := ctx.MustGet("user").(*entity.User)
	if !ok {
		ctx.Error(error.NewInternalServerError(errors.New("user couldn't get")))
	}

	ranking, err := h.rankingController.GetNearMeRanking(user)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.JSON(ctx, http.StatusOK, "", ranking, nil)
}
