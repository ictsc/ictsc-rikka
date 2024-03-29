package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	e "github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"sort"
	"time"
)

type RankingService struct {
	preRoundMode        bool
	answerLimit         time.Duration
	userGroupRepository repository.UserGroupRepository
	answerRepository    repository.AnswerRepository
}

type Rank struct {
	UserGroup   *entity.UserGroup
	Point       uint
	LastPointed time.Time
	Rank        uint
}

type problemPoint struct {
	point uint
	gotAt time.Time
}

func NewRankingService(preRoundMode bool, answerLimit int, userGroupRepository repository.UserGroupRepository, answerRepository repository.AnswerRepository) *RankingService {
	return &RankingService{
		preRoundMode:        preRoundMode,
		answerLimit:         time.Duration(answerLimit) * time.Minute,
		userGroupRepository: userGroupRepository,
		answerRepository:    answerRepository,
	}
}

func (s *RankingService) getAnswersForRanking(isFullAccess bool) ([]*entity.Answer, error) {
	answers, err := s.answerRepository.GetAll()
	if err != nil {
		return []*entity.Answer{}, err
	}

	// まだユーザに点数が公開されていない回答を除外する
	now := time.Now()
	pos := 0
	for _, answer := range answers {
		// 20分制約によって回答が見れていない場合
		if !isFullAccess && !now.After(answer.CreatedAt.Add(s.answerLimit)) {
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

func (s *RankingService) getLatestRanking(isFullAccess bool) (map[uuid.UUID]*Rank, error) {
	userGroups, err := s.userGroupRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// ユーザの順位を計算するために使うデータ構造を初期化する
	rankTable := make(map[uuid.UUID]*Rank)
	if !isFullAccess && s.preRoundMode {
		return nil, e.NewForbiddenError("preRoundMode is enabled")
	}

	for _, userGroup := range userGroups {
		// ICTSCユーザグループを除外する
		if userGroup.IsFullAccess {
			continue
		}

		rankTable[userGroup.ID] = &Rank{
			UserGroup: userGroup,
			Point:     0,
		}
	}

	// ランキング計算に利用する回答を取得する
	answers, err := s.getAnswersForRanking(isFullAccess)
	if err != nil {
		return nil, err
	}

	// 各チーム・各問題の得点を計算する
	answerTable := make(map[uuid.UUID]map[uuid.UUID]problemPoint)
	for _, answer := range answers {
		if _, ok := answerTable[answer.UserGroupID]; !ok {
			answerTable[answer.UserGroupID] = make(map[uuid.UUID]problemPoint)
		}
		point := answerTable[answer.UserGroupID][answer.ProblemID].point
		gotAt := answerTable[answer.UserGroupID][answer.ProblemID].gotAt

		// answer.Pointは、getAnswersForRankingでnullでないことが保証されているので問題ない
		if point < *answer.Point || (point == *answer.Point && answer.CreatedAt.Before(gotAt)) {
			answerTable[answer.UserGroupID][answer.ProblemID] = problemPoint{
				point: *answer.Point,
				gotAt: answer.CreatedAt,
			}
		}
	}

	// 各チームの得点を計算する
	for userGroupId, rank := range rankTable {
		for _, problemPoint := range answerTable[userGroupId] {
			rank.Point += problemPoint.point
			if problemPoint.gotAt.After(rank.LastPointed) {
				rank.LastPointed = problemPoint.gotAt
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

		if ranks[i].Point > ranks[j].Point {
			return true
		}

		if ranks[i].Point == ranks[j].Point {
			if ranks[i].LastPointed.Before(ranks[j].LastPointed) {
				return true
			}
		}

		return false
	})

	// ソートした結果を利用して順位を計算する
	ranks[0].Rank = 1
	cRank := uint(1)
	cPoint := ranks[0].Point
	cLastPointed := ranks[0].LastPointed
	for _, rank := range ranks[1:] {
		if cPoint > rank.Point || cLastPointed.Before(rank.LastPointed) {
			cRank++
		}

		cPoint = rank.Point
		cLastPointed = rank.LastPointed
		rank.Rank = cRank
	}

	return rankTable, nil
}

func (s *RankingService) getRanking(isFullAccess bool) (map[uuid.UUID]*Rank, error) {
	return s.getLatestRanking(isFullAccess)
}

func (s *RankingService) table2slice(table map[uuid.UUID]*Rank) []*Rank {
	ranks := make([]*Rank, 0, len(table))
	for _, rank := range table {
		ranks = append(ranks, rank)
	}

	sort.SliceStable(ranks, func(i, j int) bool {
		return ranks[i].Rank < ranks[j].Rank || (ranks[i].Rank == ranks[j].Rank && ranks[i].UserGroup.Name < ranks[j].UserGroup.Name)
	})

	return ranks
}

func (s *RankingService) GetRanking(group *entity.UserGroup) ([]*Rank, error) {
	var rankTable map[uuid.UUID]*Rank
	if !group.IsFullAccess && s.preRoundMode {
		return nil, e.NewForbiddenError("preRoundMode is enabled")
	}

	rankTable, err := s.getRanking(true)
	if err != nil {
		return nil, err
	}

	return s.table2slice(rankTable), nil
}

func (s *RankingService) GetTopRanking(group *entity.UserGroup) ([]*Rank, error) {
	var rankTable map[uuid.UUID]*Rank
	if !group.IsFullAccess && s.preRoundMode {
		return nil, e.NewForbiddenError("preRoundMode is enabled")
	}

	rankTable, err := s.getRanking(false)
	if err != nil {
		return nil, err
	}

	ranks := s.table2slice(rankTable)

	for i := range ranks {
		if ranks[i].Rank > 5 {
			return ranks[:i-1], nil
		}
	}

	return ranks, nil
}
