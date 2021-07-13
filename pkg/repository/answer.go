package repository

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
)

type AnswerRepository interface {
	Create(answer *entity.Answer) (*entity.Answer, error)
	GetAll() ([]*entity.Answer, error)
	FindByID(id uuid.UUID) (*entity.Answer, error)
	FindByProblem(probid uuid.UUID, teamid *uuid.UUID) ([]*entity.Answer, error)
	FindByTeam(id uuid.UUID) ([]*entity.Answer, error)
	FindByProblemAndTeam(probid uuid.UUID,teamid uuid.UUID) ([]*entity.Answer, error)
	Update(answer *entity.Answer) (*entity.Answer, error)
	Delete(answer *entity.Answer) error
}

