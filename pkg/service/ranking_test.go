package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type MockUserGroupRepository struct{}
type MockAnswerRepository struct{}

var _ repository.UserGroupRepository = (*MockUserGroupRepository)(nil)
var _ repository.AnswerRepository = (*MockAnswerRepository)(nil)

var teams []uuid.UUID
var problems []uuid.UUID
var times []time.Time

func init() {
	teams = []uuid.UUID{}
	problems = []uuid.UUID{}

	for i := 0; i < 6; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		teams = append(teams, uuid)
	}

	for i := 0; i < 3; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		problems = append(problems, uuid)
	}

	times = []time.Time{
		time.Now().Add(-70 * time.Minute),
		time.Now().Add(-60 * time.Minute),
		time.Now().Add(-50 * time.Minute),
		time.Now().Add(-40 * time.Minute),
		time.Now().Add(-30 * time.Minute),
		time.Now().Add(-20 * time.Minute),
		time.Now().Add(-10 * time.Minute),
		time.Now(),
	}

}

// Testで利用しないメソッドの定義
func (r *MockUserGroupRepository) Create(userGroup *entity.UserGroup) (*entity.UserGroup, error) {
	return nil, nil
}
func (r *MockUserGroupRepository) FindByID(id uuid.UUID) (*entity.UserGroup, error)  { return nil, nil }
func (r *MockUserGroupRepository) FindByName(name string) (*entity.UserGroup, error) { return nil, nil }
func (r *MockAnswerRepository) Create(answer *entity.Answer) (*entity.Answer, error) { return nil, nil }
func (r *MockAnswerRepository) FindByID(id uuid.UUID) (*entity.Answer, error)        { return nil, nil }
func (r *MockAnswerRepository) FindByProblem(probid uuid.UUID, groupid *uuid.UUID) ([]*entity.Answer, error) {
	return nil, nil
}
func (r *MockAnswerRepository) FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error) {
	return nil, nil
}
func (r *MockAnswerRepository) FindByProblemAndUserGroup(probid uuid.UUID, groupid uuid.UUID) ([]*entity.Answer, error) {
	return nil, nil
}
func (r *MockAnswerRepository) Update(answer *entity.Answer) (*entity.Answer, error) { return nil, nil }
func (r *MockAnswerRepository) Delete(answer *entity.Answer) error                   { return nil }

func (r *MockUserGroupRepository) GetAll() ([]*entity.UserGroup, error) {
	return []*entity.UserGroup{
		{
			Base: entity.Base{
				ID: teams[0],
			},
			Name: "team1",
		},
		{
			Base: entity.Base{
				ID: teams[1],
			},
			Name: "team2",
		},
		{
			Base: entity.Base{
				ID: teams[2],
			},
			Name: "team3",
		},
		{
			Base: entity.Base{
				ID: teams[3],
			},
			Name: "team4",
		},
		{
			Base: entity.Base{
				ID: teams[4],
			},
			Name: "team5",
		},
		{
			Base: entity.Base{
				ID: teams[5],
			},
			Name: "team6",
		},
	}, nil
}

func (r *MockAnswerRepository) GetAll() ([]*entity.Answer, error) {
	var p0 uint = 0
	var p1 uint = 100
	var p2 uint = 200
	var p3 uint = 300

	// テストシナリオ
	// - team1, prob1, 70分前, 100点
	// - team1, prob1, 60分前, 200点
	// - team1, prob1, 50分前, 100点
	// - team2, prob2, 70分前, 300点
	// - team3, prob2, 60分前, 200点
	// - team3, prob3, 10分前, 200点 (10分前なので考慮外)
	// - team4, prob1, 70分前, 100点
	// - team4, prob2, 60分前, 100点
	// - team4, prob3, 50分前, 100点
	// - team4, prob3, 30分前, nil (点数がnilなので考慮外)
	// - team5, prob1, 30分前, 0点
	// 順位
	// 1. team2 300点 (60分前)
	// 2. team4 300点 (50分前)
	// 3. team1 200点 (60分前)
	// 3. team3 200点 (60分前)
	// 4. team5 0点
	// 4. team6 0点

	return []*entity.Answer{
		// team1, prob1, 70分前, 100点
		{
			Base: entity.Base{
				CreatedAt: times[0],
			},
			Point:       &p1,
			UserGroupID: teams[0],
			ProblemID:   problems[0],
		},
		// team1, prob1, 60分前, 200点
		{
			Base: entity.Base{
				CreatedAt: times[1],
			},
			Point:       &p2,
			UserGroupID: teams[0],
			ProblemID:   problems[0],
		},
		// team1, prob1, 50分前, 100点
		{
			Base: entity.Base{
				CreatedAt: times[2],
			},
			Point:       &p1,
			UserGroupID: teams[0],
			ProblemID:   problems[0],
		},
		// team2, prob2, 70分前, 300点
		{
			Base: entity.Base{
				CreatedAt: times[0],
			},
			Point:       &p3,
			UserGroupID: teams[1],
			ProblemID:   problems[1],
		},
		// team3, prob2, 60分前, 200点
		{
			Base: entity.Base{
				CreatedAt: times[1],
			},
			Point:       &p2,
			UserGroupID: teams[2],
			ProblemID:   problems[1],
		},
		// team3, prob3, 10分前, 200点
		{
			Base: entity.Base{
				CreatedAt: times[6],
			},
			Point:       &p2,
			UserGroupID: teams[2],
			ProblemID:   problems[2],
		},
		// team4, prob1, 70分前, 100点
		{
			Base: entity.Base{
				CreatedAt: times[0],
			},
			Point:       &p1,
			UserGroupID: teams[3],
			ProblemID:   problems[0],
		},
		// team4, prob2, 60分前, 100点
		{
			Base: entity.Base{
				CreatedAt: times[1],
			},
			Point:       &p1,
			UserGroupID: teams[3],
			ProblemID:   problems[1],
		},
		// team4, prob3, 50分前, 100点
		{
			Base: entity.Base{
				CreatedAt: times[2],
			},
			Point:       &p1,
			UserGroupID: teams[3],
			ProblemID:   problems[2],
		},
		// team4, prob3, 30分前, nil
		{
			Base: entity.Base{
				CreatedAt: times[4],
			},
			Point:       nil,
			UserGroupID: teams[3],
			ProblemID:   problems[2],
		},
		// team5, prob1, 30分前, 0点
		{
			Base: entity.Base{
				CreatedAt: times[2],
			},
			Point:       &p0,
			UserGroupID: teams[4],
			ProblemID:   problems[0],
		},
	}, nil
}

func TestRanking(t *testing.T) {
	s := NewRankingService(&MockUserGroupRepository{}, &MockAnswerRepository{})

	rankTable, err := s.getLatestRanking()
	if err != nil {
		t.Error(err)
	}

	// 順位
	// 1. team2 300点 (60分前)
	// 2. team4 300点 (50分前)
	// 3. team1 200点 (60分前)
	// 3. team3 200点 (60分前)
	// 4. team5 0点
	// 4. team6 0点
	{
		rank := rankTable[teams[0]]
		if rank.Rank != 3 {
			t.Errorf("team1's rank expect 3, actual %d", rank.Rank)
		}
		if rank.Point != 200 {
			t.Errorf("team1's point expect 200, actual %d", rank.Point)
		}
	}

	{
		rank := rankTable[teams[1]]
		if rank.Rank != 1 {
			t.Errorf("team2's rank expect 1, actual %d", rank.Rank)
		}
		if rank.Point != 300 {
			t.Errorf("team2's point expect 300, actual %d", rank.Point)
		}
	}

	{
		rank := rankTable[teams[2]]
		if rank.Rank != 3 {
			t.Errorf("team3's rank expect 3, actual %d", rank.Rank)
		}
		if rank.Point != 200 {
			t.Errorf("team3's point expect 200, actual %d", rank.Point)
		}
	}

	{
		rank := rankTable[teams[3]]
		if rank.Rank != 2 {
			t.Errorf("team4's rank expect 2, actual %d", rank.Rank)
		}
		if rank.Point != 300 {
			t.Errorf("team4's point expect 300, actual %d", rank.Point)
		}
	}

	{
		rank := rankTable[teams[4]]
		if rank.Rank != 4 {
			t.Errorf("team5's rank expect 4, actual %d", rank.Rank)
		}
		if rank.Point != 0 {
			t.Errorf("team5's point expect 0, actual %d", rank.Point)
		}
	}

	{
		rank := rankTable[teams[5]]
		if rank.Rank != 4 {
			t.Errorf("team6's rank expect 4, actual %d", rank.Rank)
		}
		if rank.Point != 0 {
			t.Errorf("team6's point expect 0, actual %d", rank.Point)
		}
	}
}
