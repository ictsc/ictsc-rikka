package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type AnswerService struct {
	userRepo    repository.UserRepository
	answerRepo  repository.AnswerRepository
	problemRepo repository.ProblemRepository
}

type CreateAnswerRequest struct {
	Point     uint
	Body      string
	Group     uuid.UUID
	ProblemID uuid.UUID
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
	ans := &entity.Answer{
		UserGroupID: req.Group,
		Point:       0,
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

	ans.Point = req.Point

	return s.answerRepo.Update(ans)
}

func (s *AnswerService) Delete(id uuid.UUID) error {
	ans, err := s.answerRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.answerRepo.Delete(ans)
}
