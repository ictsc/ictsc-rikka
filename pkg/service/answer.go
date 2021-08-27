package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	e "github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type AnswerService struct {
	answerLimit time.Duration
	userRepo    repository.UserRepository
	answerRepo  repository.AnswerRepository
	problemRepo repository.ProblemRepository
}

type CreateAnswerRequest struct {
	UserGroup *entity.UserGroup
	Body      string
	ProblemID uuid.UUID
}

type UpdateAnswerRequest struct {
	Point uint
}

func NewAnswerService(answerLimit int, userRepo repository.UserRepository, answerRepo repository.AnswerRepository, problemRepo repository.ProblemRepository) *AnswerService {
	return &AnswerService{
		answerLimit: time.Duration(answerLimit) * time.Minute,
		userRepo:    userRepo,
		answerRepo:  answerRepo,
		problemRepo: problemRepo,
	}
}

func (s *AnswerService) Create(req *CreateAnswerRequest) (*entity.Answer, error) {
	lastAnswered := time.Time{}
	pastAnswers, err := s.answerRepo.FindByUserGroup(req.UserGroup.ID)
	if err != nil {
		return nil, err
	}

	for _, answer := range pastAnswers {
		if answer.CreatedAt.After(lastAnswered) && answer.ProblemID == req.ProblemID {
			lastAnswered = answer.CreatedAt
		}
	}

	if !time.Now().After(lastAnswered.Add(s.answerLimit)) {
		return nil, e.NewForbiddenError("couldn't submit answer if you submit answer within last 20 minutes")
	}

	ans := &entity.Answer{
		UserGroupID: req.UserGroup.ID,
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

func (s *AnswerService) FindByID(group *entity.UserGroup, id uuid.UUID) (*entity.Answer, error) {
	ans, err := s.answerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !group.IsFullAccess && !time.Now().After(ans.CreatedAt.Add(s.answerLimit)) {
		ans.Point = nil
	}

	return ans, nil
}

// userGroupID is optional
func (s *AnswerService) FindByProblem(group *entity.UserGroup, probid uuid.UUID, userGroupID *uuid.UUID) ([]*entity.Answer, error) {
	if !group.IsFullAccess && userGroupID != nil && group.ID != *userGroupID {
		return nil, e.NewForbiddenError("you cannot fetch other group's answers")
	}

	answers, err := s.answerRepo.FindByProblem(probid, userGroupID)
	if err != nil {
		return nil, err
	}

	if !group.IsFullAccess {
		now := time.Now()
		for _, ans := range answers {
			if !now.After(ans.CreatedAt.Add(s.answerLimit)) {
				ans.Point = nil
			}
		}
	}

	return answers, nil
}

func (s *AnswerService) FindByUserGroup(id uuid.UUID) ([]*entity.Answer, error) {
	return s.answerRepo.FindByUserGroup(id)
}

// userGroupID is require
func (s *AnswerService) FindByProblemAndUserGroup(group *entity.UserGroup, probid uuid.UUID, userGroupID uuid.UUID) ([]*entity.Answer, error) {
	if !group.IsFullAccess && group.ID != userGroupID {
		return nil, e.NewForbiddenError("you cannot fetch other group's answers")
	}

	answers, err := s.answerRepo.FindByProblemAndUserGroup(probid, userGroupID)
	if err != nil {
		return nil, err
	}

	if !group.IsFullAccess {
		now := time.Now()
		for _, ans := range answers {
			if !now.After(ans.CreatedAt.Add(s.answerLimit)) {
				ans.Point = nil
			}
		}
	}

	return answers, nil
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
	if problem == nil {
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
