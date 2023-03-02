package rc

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type NoticeWithSyncTimeRepository struct {
	rc *redis.Client
}

func NewNoticeWithSyncTimeRepository(rc *redis.Client) *NoticeWithSyncTimeRepository {
	return &NoticeWithSyncTimeRepository{rc: rc}
}

func (r *NoticeWithSyncTimeRepository) Set(ctx context.Context, path string, noticeWithInfo entity.NoticeWithSyncTime) error {
	jsonBytes, err := json.Marshal(noticeWithInfo)
	if err != nil {
		return err
	}

	err = r.rc.Set(ctx, path, jsonBytes, 0).Err()
	return err
}

func (r *NoticeWithSyncTimeRepository) Get(ctx context.Context, path string) (*entity.NoticeWithSyncTime, error) {
	noticeWithInfo, err := r.rc.Get(ctx, path).Result()
	if err != nil {
		return nil, err
	}

	noticeWithInfoEntity := entity.NoticeWithSyncTime{}
	err = json.Unmarshal([]byte(noticeWithInfo), &noticeWithInfoEntity)

	return &noticeWithInfoEntity, err
}
