package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

func (r *AnswerRepository) Create(answer *entity.Answer) (*entity.Answer, error) {
	err := r.db.Create(answer).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(answer.ID)
}

func (r *AnswerRepository) GetAll() ([]*entity.Answer, error) {
	answers := make([]*entity.Answer, 0)
	err := r.db.Find(answers).Error
	return answers, err
}

func (r *AnswerRepository) FindByID(id uuid.UUID) (*entity.Answer, error) {
	res := &entity.Answer{}
	err := r.db.First(res, id).Error
	return res, err
}

func (r *AnswerRepository) FindByProblem(id uuid.UUID) (*entity.Answer, error) {
	res := &entity.Answer{}
	err := r.db.Where("problem", id).First(res).Error
	return res, err
}

func (r *AnswerRepository) FindByTeam(id uuid.UUID) (*entity.Answer, error) {
	res := &entity.Answer{}
	err := r.db.Where("team", id).First(res).Error
	return res, err
}

func (r *AnswerRepository) FindByProblemAndTeam(problemid uuid.UUID,teamid uuid.UUID) (*entity.Answer, error) {
	res := &entity.Answer{}
	err := r.db.Where("team", problemid).Where("problem",teamid).First(res).Error
	return res, err
}

func (r *AnswerRepository) Update(answer *entity.Answer) (*entity.Answer, error) {
	err := r.db.Save(answer).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(answer.ID)
}

func (r *AnswerRepository) Delete(answer *entity.Answer) error {
	return r.db.Delete(answer, answer.ID).Error
}
