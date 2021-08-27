package mariadb

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
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

var _ repository.AnswerRepository = (*AnswerRepository)(nil)

func (r *AnswerRepository) Create(answer *entity.Answer) (*entity.Answer, error) {
	err := r.db.Create(answer).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(answer.ID)
}

func (r *AnswerRepository) GetAll() ([]*entity.Answer, error) {
	answers := make([]*entity.Answer, 0)
	err := r.db.Preload("UserGroup").Find(&answers).Error
	return answers, err
}

func (r *AnswerRepository) FindByID(id uuid.UUID) (*entity.Answer, error) {
	res := &entity.Answer{}
	err := r.db.Preload("UserGroup").First(res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *AnswerRepository) FindByProblem(probid uuid.UUID, groupid *uuid.UUID) ([]*entity.Answer, error) {
	res := []*entity.Answer{}
	if groupid != nil {
		err := r.db.Preload("UserGroup").Where("problem_id", probid).Where("user_group_id", groupid).Find(&res).Error
		return res, err
	} else {
		err := r.db.Preload("UserGroup").Where("problem_id", probid).Find(&res).Error
		return res, err
	}
}

func (r *AnswerRepository) FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error) {
	res := []*entity.Answer{}
	err := r.db.Preload("UserGroup").Where("user_group_id", id).Find(&res).Error
	return res, err
}

func (r *AnswerRepository) FindByProblemAndUserGroup(problemid uuid.UUID, groupid uuid.UUID) ([]*entity.Answer, error) {
	res := []*entity.Answer{}
	err := r.db.Preload("UserGroup").Where("problem_id", problemid).Where("user_group_id", groupid).Find(&res).Error
	return res, err
}

func (r *AnswerRepository) Update(answer *entity.Answer) (*entity.Answer, error) {
	err := r.db.Model(&entity.Answer{}).Where("id", answer.ID).Update("point", answer.Point).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(answer.ID)
}

func (r *AnswerRepository) Delete(answer *entity.Answer) error {
	return r.db.Delete(answer, answer.ID).Error
}
