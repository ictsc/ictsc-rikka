package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type AnswerRepository interface {
	Create(answer *entity.Answer) (*entity.Answer, error)
	GetAll() ([]*entity.Answer, error)
	FindByID(id uuid.UUID) (*entity.Answer, error)
	FindByProblem(probid uuid.UUID, groupid *uuid.UUID) ([]*entity.Answer, error)
	FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error)
	FindByProblemAndUserGroup(probid uuid.UUID, groupid uuid.UUID) ([]*entity.Answer, error)
	Update(answer *entity.Answer) (*entity.Answer, error)
	Delete(answer *entity.Answer) error
}
