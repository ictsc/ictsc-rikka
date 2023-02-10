package controller

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type RankingController struct {
	rankingService *service.RankingService
}

func NewRankingController(rankingService *service.RankingService) *RankingController {
	return &RankingController{
		rankingService: rankingService,
	}
}

type RankingEntry struct {
	UserGroupID uuid.UUID         `json:"user_group_id"`
	UserGroup   *entity.UserGroup `json:"user_group"`
	Rank        uint              `json:"rank"`
	Point       uint              `json:"point"`
}

type Ranking struct {
	Ranking []RankingEntry `json:"ranking"`
}

func (c *RankingController) entity2response(ranks []*service.Rank) *Ranking {
	ranks_resp := make([]RankingEntry, 0, len(ranks))
	for _, rank := range ranks {
		ranks_resp = append(ranks_resp, RankingEntry{
			UserGroupID: rank.UserGroup.ID,
			UserGroup:   rank.UserGroup,
			Rank:        rank.Rank,
			Point:       rank.Point,
		})
	}
	return &Ranking{Ranking: ranks_resp}
}

func (c *RankingController) GetRanking() (*Ranking, error) {
	ranking, err := c.rankingService.GetRanking()
	if err != nil {
		return nil, err
	}

	return c.entity2response(ranking), nil
}

func (c *RankingController) GetTopRanking() (*Ranking, error) {
	ranking, err := c.rankingService.GetTopRanking()
	if err != nil {
		return nil, err
	}

	return c.entity2response(ranking), nil
}
