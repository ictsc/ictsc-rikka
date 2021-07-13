package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type AnswerService struct {
	userRepo      repository.UserRepository
	answerRepo repository.AnswerRepository
	problemRepo repository.ProblemRepository
}

type CreateAnswerRequest struct {
	Point             uint
	Body           string
	ProblemID *uuid.UUID
}

type UpdateAnswerRequest struct {
	Point             uint
	Body           string
	ProblemID *uuid.UUID
}

func NewAnswerService(userRepo repository.UserRepository, answerRepo repository.AnswerRepository, problemRepo repository.ProblemRepository) *AnswerService {
	return &AnswerService{
		userRepo:      userRepo,
		answerRepo: answerRepo,
		problemRepo: problemRepo,
	}
}

func (s *AnswerService) Create(req *CreateAnswerRequest) (*entity.Answer, error) {
	ans := &entity.Answer{
		Point: req.Point,
		Body: req.Body,
		ProblemID: req.ProblemID,
	}

	if ans.ProblemID != nil {
		return nil, errors.New("problem id is empty")
	}
	problem, err := s.problemRepo.FindByID(*req.ProblemID)
	if err != nil {
		return nil, err
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

func (s *AnswerService) FindByProblem(id uuid.UUID) (*entity.Answer, error) {
	return s.answerRepo.FindByProblem(id)
}

func (s *AnswerService) FindByTeam(id uuid.UUID) (*entity.Answer, error) {
	return s.answerRepo.FindByTeam(id)
}

func (s *AnswerService) FindByProblemAndTeam(probid uuid.UUID, teamid uuid.UUID) (*entity.Answer, error) {
	return s.answerRepo.FindByProblemAndTeam(probid,teamid)
}

func (s *AnswerService) Update(id uuid.UUID, req *UpdateAnswerRequest) (*entity.Answer, error) {
	ans, err := s.answerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ans == nil {
		return nil, errors.New("answer not found")
	}
	if req.ProblemID != nil {
		return nil, errors.New("problem id can not be changed")
	}

	ans.Point = req.Point
	ans.Body = req.Body

	return s.answerRepo.Update(ans)
}

func (s *AnswerService) Delete(id uuid.UUID) error {
	ans, err := s.answerRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.answerRepo.Delete(ans)
}
