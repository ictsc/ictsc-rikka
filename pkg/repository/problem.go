package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemRepository interface {
	Create(problem *entity.Problem) (*entity.Problem, error)
	GetAll() ([]*entity.Problem, error)
	FindByID(id uuid.UUID) (*entity.Problem, error)
	FindByCode(code string) (*entity.Problem, error)
	Update(problem *entity.Problem) (*entity.Problem, error)
	Delete(problem *entity.Problem) error
}
