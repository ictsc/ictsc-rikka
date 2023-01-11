package rc

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemWithInfoRepository struct {
	rc *redis.Client
}

func NewProblemWithInfoRepository(rc *redis.Client) *ProblemWithInfoRepository {
	return &ProblemWithInfoRepository{rc: rc}
}

func (r *ProblemWithInfoRepository) Set(ctx context.Context, path string, problemWithInfo entity.ProblemWithInfo) error {
	jsonBytes, err := json.Marshal(problemWithInfo)
	if err != nil {
		return err
	}

	err = r.rc.Set(ctx, path, jsonBytes, 0).Err()
	return err
}

func (r *ProblemWithInfoRepository) Get(ctx context.Context, path string) (*entity.ProblemWithInfo, error) {
	problemWithInfo, err := r.rc.Get(ctx, path).Result()
	if err != nil {
		return nil, err
	}

	problemWithInfoEntity := entity.ProblemWithInfo{}
	err = json.Unmarshal([]byte(problemWithInfo), &problemWithInfoEntity)

	return &problemWithInfoEntity, err
}
