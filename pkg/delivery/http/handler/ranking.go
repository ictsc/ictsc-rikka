package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/response"
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
	{
		authed := route.Group("")

		authed.Use(middleware.Auth(userRepo))

		authed.GET("", handler.GetRanking)
	}
}

func (h *RankingHandler) GetRanking(ctx *gin.Context) {
	ranking, err := h.rankingController.GetRanking()
	if err != nil {
		ctx.Error(err)
		return
	}

	response.JSON(ctx, http.StatusOK, "", ranking, nil)
}
