package service

import (
<<<<<<< HEAD
	"fmt"
=======
	"time"
>>>>>>> implement GET /usergroups

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	e "github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type AnswerService struct {
	userRepo    repository.UserRepository
	answerRepo  repository.AnswerRepository
	problemRepo repository.ProblemRepository
}

type CreateAnswerRequest struct {
	Point       uint
	Body        string
	UserGroupID uuid.UUID
	ProblemID   uuid.UUID
}

type UpdateAnswerRequest struct {
	Point uint
}

func NewAnswerService(userRepo repository.UserRepository, answerRepo repository.AnswerRepository, problemRepo repository.ProblemRepository) *AnswerService {
	return &AnswerService{
		userRepo:    userRepo,
		answerRepo:  answerRepo,
		problemRepo: problemRepo,
	}
}

func (s *AnswerService) Create(req *CreateAnswerRequest) (*entity.Answer, error) {
	lastAnswered := time.Time{}
	pastAnswers, err := s.answerRepo.FindByUserGroup(req.UserGroupID)
	if err != nil {
		return nil, err
	}

	for _, answer := range pastAnswers {
		if answer.CreatedAt.After(lastAnswered) {
			lastAnswered = answer.CreatedAt
		}
	}

	if !time.Now().After(lastAnswered.Add(1 * time.Minute)) {
		return nil, e.NewForbiddenError("couldn't submit answer if you submit answer within last 20 minutes")
	}

	ans := &entity.Answer{
		UserGroupID: req.UserGroupID,
		Point:       nil,
		Body:        req.Body,
		ProblemID:   req.ProblemID,
	}

	problem, err := s.problemRepo.FindByID(req.ProblemID)
	if err != nil {
		return nil, errors.New("problem id is invalid")
	}
	if problem == nil {
		return nil, errors.New("problem id is invalid")
	}

	return s.answerRepo.Create(ans)
}

func (s *AnswerService) GetAll() ([]*entity.Answer, error) {
	return s.answerRepo.GetAll()
}

func (s *AnswerService) FindByID(id uuid.UUID) (*entity.Answer, error) {
	return s.answerRepo.FindByID(id)
}

// userGroupID is optional
func (s *AnswerService) FindByProblem(probid uuid.UUID, userGroupID *uuid.UUID) ([]*entity.Answer, error) {
	return s.answerRepo.FindByProblem(probid, userGroupID)
}

func (s *AnswerService) FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error) {
	return s.answerRepo.FindByUserGroup(id)
}

// userGroupID is require
func (s *AnswerService) FindByProblemAndUserGroup(probid uuid.UUID, userGroupID uuid.UUID) ([]*entity.Answer, error) {
	return s.answerRepo.FindByProblemAndUserGroup(probid, userGroupID)
}

func (s *AnswerService) Update(id uuid.UUID, req *UpdateAnswerRequest) (*entity.Answer, error) {
	ans, err := s.answerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ans == nil {
		return nil, errors.New("answer not found")
	}

	problem, err := s.problemRepo.FindByID(ans.ProblemID)
	if err != nil {
		return nil, e.NewInternalServerError(err)
	}
	if ans == nil {
		return nil, e.NewInternalServerError(fmt.Errorf("problem %s bound answer %s is not found", ans.ProblemID, ans.ID))
	}

	if !(req.Point <= problem.Point) {
		return nil, e.NewBadRequestError("invalid point")
	}

	ans.Point = &req.Point

	return s.answerRepo.Update(ans)
}

func (s *AnswerService) Delete(id uuid.UUID) error {
	ans, err := s.answerRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.answerRepo.Delete(ans)
}
