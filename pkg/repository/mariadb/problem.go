package mariadb

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type ProblemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) *ProblemRepository {
	return &ProblemRepository{
		db: db,
	}
}

func (r *ProblemRepository) Create(problem *entity.Problem) (*entity.Problem, error) {
	err := r.db.Create(problem).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(problem.ID)
}

func (r *ProblemRepository) GetAll() ([]*entity.Problem, error) {
	problems := make([]*entity.Problem, 0)
	err := r.db.Find(&problems).Error
	return problems, err
}

func (r *ProblemRepository) FindByID(id uuid.UUID) (*entity.Problem, error) {
	res := &entity.Problem{}
	err := r.db.First(res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *ProblemRepository) FindByCode(code string) (*entity.Problem, error) {
	res := &entity.Problem{}
	err := r.db.Where("code", code).First(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *ProblemRepository) Update(problem *entity.Problem) (*entity.Problem, error) {
	err := r.db.Save(problem).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(problem.ID)
}

func (r *ProblemRepository) Delete(problem *entity.Problem) error {
	return r.db.Delete(problem, problem.ID).Error
}
