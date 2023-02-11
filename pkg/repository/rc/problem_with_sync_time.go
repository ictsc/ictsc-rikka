package rc

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemWithSyncTimeRepository struct {
	rc *redis.Client
}

func NewProblemWithSyncTimeRepository(rc *redis.Client) *ProblemWithSyncTimeRepository {
	return &ProblemWithSyncTimeRepository{rc: rc}
}

func (r *ProblemWithSyncTimeRepository) Set(ctx context.Context, path string, problemWithInfo entity.ProblemWithSyncTime) error {
	jsonBytes, err := json.Marshal(problemWithInfo)
	if err != nil {
		return err
	}

	err = r.rc.Set(ctx, path, jsonBytes, 0).Err()
	return err
}

func (r *ProblemWithSyncTimeRepository) Get(ctx context.Context, path string) (*entity.ProblemWithSyncTime, error) {
	problemWithInfo, err := r.rc.Get(ctx, path).Result()
	if err != nil {
		return nil, err
	}

	problemWithInfoEntity := entity.ProblemWithSyncTime{}
	err = json.Unmarshal([]byte(problemWithInfo), &problemWithInfoEntity)

	return &problemWithInfoEntity, err
}
