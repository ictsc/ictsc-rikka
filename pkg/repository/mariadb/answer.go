package mariadb

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type AnswerRepository struct {
	*db
}

func NewAnswerRepository(config *MariaDBConfig) *AnswerRepository {
	return &AnswerRepository{
		db: newDB(config),
	}
}

var _ repository.AnswerRepository = (*AnswerRepository)(nil)

func (r *AnswerRepository) Create(answer *entity.Answer) (*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = db.Create(answer).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(answer.ID)
}

func (r *AnswerRepository) GetAll() ([]*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	answers := make([]*entity.Answer, 0)
	err = db.Find(answers).Error
	return answers, err
}

func (r *AnswerRepository) FindByID(id uuid.UUID) (*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := &entity.Answer{}
	err = db.First(res, id).Error
	return res, err
}

func (r *AnswerRepository) FindByProblem(probid uuid.UUID, groupid *uuid.UUID) ([]*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := []*entity.Answer{}
	if groupid != nil {
		err = db.Where("problem_id", probid).Where("user_group_id", groupid).Find(res).Error
	} else {
		err = db.Where("problem_id", probid).Find(&res).Error
	}
	return res, err
}

func (r *AnswerRepository) FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := []*entity.Answer{}
	err = db.Where("user_group_id", id).Find(&res).Error
	return res, err
}

func (r *AnswerRepository) FindByProblemAndUserGroup(problemid uuid.UUID, groupid uuid.UUID) ([]*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := []*entity.Answer{}
	err = db.Where("problem_id", problemid).Where("user_group_id", groupid).Find(&res).Error
	return res, err
}

func (r *AnswerRepository) Update(answer *entity.Answer) (*entity.Answer, error) {
	db, conn, err := r.db.init()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = db.Save(answer).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(answer.ID)
}

func (r *AnswerRepository) Delete(answer *entity.Answer) error {
	db, conn, err := r.db.init()
	if err != nil {
		return err
	}
	defer conn.Close()

	return db.Delete(answer, answer.ID).Error
}
