package repository

import (
	"context"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemWithInfoRepository interface {
	Set(context.Context, string, entity.ProblemWithInfo) error
	Get(context.Context, string) (*entity.ProblemWithInfo, error)
}
