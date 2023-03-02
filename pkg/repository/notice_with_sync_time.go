package repository

import (
	"context"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type NoticeWithSyncTimeRepository interface {
	Set(context.Context, string, entity.NoticeWithSyncTime) error
	Get(context.Context, string) (*entity.NoticeWithSyncTime, error)
}
