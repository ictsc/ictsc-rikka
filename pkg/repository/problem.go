//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemRepository interface {
	Create(problem *entity.Problem) (*entity.Problem, error)
	GetAll() ([]*entity.Problem, error)
	GetProblemsWithIsAnsweredByUserGroup(uuid.UUID) ([]*entity.ProblemWithIsAnswered, error)
	FindByID(id uuid.UUID) (*entity.Problem, error)
	FindByCode(code string) (*entity.Problem, error)
	Update(problem *entity.Problem, skipUpdatedAt bool) (*entity.Problem, error)
	Delete(problem *entity.Problem) error
}
