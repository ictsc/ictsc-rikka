package handler

import (
	"fmt"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
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
		group := route.Group("")
		authed := group.Group("")
		authed.Use(middleware.Auth(userRepo))

		authed.GET("", handler.GetRanking)
		authed.GET("/top", handler.GetTopRanking)
	}
}

func (h *RankingHandler) GetRanking(ctx *gin.Context) {
	fmt.Println("テスト1")
	group := ctx.MustGet("group").(*entity.UserGroup)
	fmt.Println(group)
	ranking, err := h.rankingController.GetRanking(group)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.JSON(ctx, http.StatusOK, "", ranking, nil)
}

func (h *RankingHandler) GetTopRanking(ctx *gin.Context) {
	group := ctx.MustGet("group").(*entity.UserGroup)
	ranking, err := h.rankingController.GetTopRanking(group)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.JSON(ctx, http.StatusOK, "", ranking, nil)
}
