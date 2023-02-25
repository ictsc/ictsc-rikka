package repository

import (
	"context"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemWithSyncTimeRepository interface {
	Set(context.Context, string, entity.ProblemWithSyncTime) error
	Get(context.Context, string) (*entity.ProblemWithSyncTime, error)
}
