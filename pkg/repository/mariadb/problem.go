package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type ProblemRepository struct {
	*db
}

func NewProblemRepository(config *MariaDBConfig) *ProblemRepository {
	return &ProblemRepository{
		db: newDB(config),
	}
}

func (r *ProblemRepository) Create(problem *entity.Problem) (*entity.Problem, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = db.Create(problem).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(problem.ID)
}

func (r *ProblemRepository) GetAll() ([]*entity.Problem, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	problems := make([]*entity.Problem, 0)
	err = db.Find(&problems).Error
	return problems, err
}

func (r *ProblemRepository) FindByID(id uuid.UUID) (*entity.Problem, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.Problem{}
	err = db.First(res, id).Error
	return res, err
}

func (r *ProblemRepository) FindByCode(code string) (*entity.Problem, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.Problem{}
	err = db.Where("code", code).First(res).Error
	return res, err
}

func (r *ProblemRepository) Update(problem *entity.Problem) (*entity.Problem, error) {
	db, conn, err := r.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = db.Save(problem).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(problem.ID)
}

func (r *ProblemRepository) Delete(problem *entity.Problem) error {
	db, conn, err := r.init()
	if err != nil {
		return err
	}
	defer conn.Close()

	return db.Delete(problem, problem.ID).Error
}
