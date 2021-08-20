package service

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	e "github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type RankingService struct {
	userGroupRepository repository.UserGroupRepository
	answerRepository    repository.AnswerRepository
}

type Rank struct {
	userGroup   *entity.UserGroup
	point       uint
	lastPointed time.Time
	rank        uint
}

type problemPoint struct {
	point uint
	gotAt time.Time
}

func NewRankingService(userGroupRepository repository.UserGroupRepository, answerRepository repository.AnswerRepository) *RankingService {
	return &RankingService{
		userGroupRepository: userGroupRepository,
		answerRepository:    answerRepository,
	}
}

func (s *RankingService) getAnswersForRanking() ([]*entity.Answer, error) {
	answers, err := s.answerRepository.GetAll()
	if err != nil {
		return []*entity.Answer{}, err
	}

	// まだユーザに点数が公開されていない回答を除外する
	pos := 0
	for _, answer := range answers {
		// 20分制約によって回答が見れていない場合
		if !time.Now().After(answer.CreatedAt.Add(20 * time.Minute)) {
			continue
		}

		// 採点がされていない
		if answer.Point == nil {
			continue
		}

		answers[pos] = answer
		pos++
	}

	return answers[:pos], nil
}

func (s *RankingService) getLatestRanking() (map[uuid.UUID]*Rank, error) {
	userGroups, err := s.userGroupRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// ユーザの順位を計算するために使うデータ構造を初期化する
	rankTable := make(map[uuid.UUID]*Rank)
	for _, userGroup := range userGroups {
		// ICTSCユーザグループを除外する
		if userGroup.IsFullAccess {
			continue
		}

		rankTable[userGroup.ID] = &Rank{
			userGroup: userGroup,
			point:     0,
		}
	}

	// ランキング計算に利用する回答を取得する
	answers, err := s.getAnswersForRanking()
	if err != nil {
		return nil, err
	}

	// 各チーム・各問題の得点を計算する
	answerTable := make(map[uuid.UUID]map[uuid.UUID]problemPoint)
	for _, userGroup := range userGroups {
		answerTable[userGroup.ID] = make(map[uuid.UUID]problemPoint)
	}

	for _, answer := range answers {
		point := answerTable[answer.UserGroupID][answer.ProblemID].point

		// answer.Pointは、getAnswersForRankingでnullでないことが保証されているので問題ない
		if point < *answer.Point {
			answerTable[answer.UserGroupID][answer.ProblemID] = problemPoint{
				point: *answer.Point,
				gotAt: answer.CreatedAt,
			}
		}
	}

	// 各チームの得点を計算する
	for userGroupId := range rankTable {
		rank := rankTable[userGroupId]
		for _, problemPoint := range answerTable[userGroupId] {
			rank.point += problemPoint.point
			if problemPoint.gotAt.After(rank.lastPointed) {
				rank.lastPointed = problemPoint.gotAt
			}
		}
	}

	// ソートするためにマップからスライスに変換する
	ranks := make([]*Rank, 0, len(rankTable))
	for _, rank := range rankTable {
		ranks = append(ranks, rank)
	}

	// ソートする
	sort.SliceStable(ranks, func(i, j int) bool {
		// 1. 得点順でソートする
		// 2. その点数になる最後の加点が行われた回答の投稿日時でソートする

		if ranks[i].point > ranks[j].point {
			return true
		}

		if ranks[i].point == ranks[j].point {
			if ranks[i].lastPointed.Before(ranks[j].lastPointed) {
				return true
			}
		}

		return false
	})

	// ソートした結果を利用して順位を計算する
	ranks[0].rank = 1
	cRank := uint(1)
	cPoint := ranks[0].point
	cLastPointed := ranks[0].lastPointed
	for _, rank := range ranks[1:] {
		if cPoint > rank.point || cLastPointed.Before(rank.lastPointed) {
			cRank++
		}

		cPoint = rank.point
		cLastPointed = rank.lastPointed
		rank.rank = cRank
	}

	return rankTable, nil
}

func (s *RankingService) getRanking() (map[uuid.UUID]*Rank, error) {
	return s.getLatestRanking()
}

func (s *RankingService) table2slice(table map[uuid.UUID]*Rank) []*Rank {
	ranks := make([]*Rank, 0, len(table))
	for _, rank := range table {
		ranks = append(ranks, rank)
	}

	sort.SliceStable(ranks, func(i, j int) bool {
		return ranks[i].rank < ranks[j].rank
	})

	return ranks
}

func (s *RankingService) GetTopRanking() ([]*Rank, error) {
	rankTable, err := s.getRanking()
	if err != nil {
		return nil, err
	}

	ranks := s.table2slice(rankTable)

	// 上位5チームを抽出する
	return ranks[:5], nil
}

func (s *RankingService) GetNearMeRanking(user *entity.User) ([]*Rank, error) {
	rankTable, err := s.getRanking()
	if err != nil {
		return nil, err
	}

	ranks := s.table2slice(rankTable)

	cRank, ok := rankTable[user.UserGroupID]
	if !ok {
		return nil, e.NewInternalServerError(errors.New("user group not found"))
	}

	min := cRank.rank - 1
	max := cRank.rank + 1

	pos := 0
	for _, rank := range ranks {
		if min <= rank.rank && rank.rank <= max {
			ranks[pos] = rank
			pos++
		}
	}

	return ranks[:pos], nil
}
