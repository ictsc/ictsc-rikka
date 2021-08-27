package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type ProblemService struct {
	uncheckedNearOverdueThreshold time.Duration
	uncheckedOverdueThreshold     time.Duration

	userRepo    repository.UserRepository
	problemRepo repository.ProblemRepository
	answerRepo  repository.AnswerRepository
}

type CreateProblemRequest struct {
	Code              string
	AuthorID          uuid.UUID
	Title             string
	Body              string
	Point             uint
	PreviousProblemID *uuid.UUID
	SolvedCriterion   uint
}

type UpdateProblemRequest struct {
	ID                uuid.UUID
	Code              string
	AuthorID          uuid.UUID
	Title             string
	Body              string
	Point             uint
	PreviousProblemID *uuid.UUID
	SolvedCriterion   uint
}

func NewProblemService(answerLimit int, userRepo repository.UserRepository, problemRepo repository.ProblemRepository, answerRepo repository.AnswerRepository) *ProblemService {
	return &ProblemService{
		uncheckedNearOverdueThreshold: time.Duration(answerLimit*3/4) * time.Minute,
		uncheckedOverdueThreshold:     time.Duration(answerLimit) * time.Minute,

		userRepo:    userRepo,
		problemRepo: problemRepo,
		answerRepo:  answerRepo,
	}
}

func (s *ProblemService) Create(req *CreateProblemRequest) (*entity.Problem, error) {
	prob := &entity.Problem{
		Code:              req.Code,
		AuthorID:          req.AuthorID,
		Title:             req.Title,
		Body:              req.Body,
		Point:             req.Point,
		PreviousProblemID: req.PreviousProblemID,
		SolvedCriterion:   req.SolvedCriterion,
	}

	if err := prob.Validate(); err != nil {
		return nil, err
	}

	if prob.PreviousProblemID != nil {
		prevProb, err := s.problemRepo.FindByID(*req.PreviousProblemID)
		if err != nil {
			return nil, err
		}
		if prevProb == nil {
			return nil, errors.New("previous_problem not found")
		}
	}

	return s.problemRepo.Create(prob)
}

func (s *ProblemService) GetAll() ([]*entity.Problem, error) {
	return s.problemRepo.GetAll()
}

func (s *ProblemService) GetAllWithAnswerInformation() ([]*entity.ProblemWithAnswerInformation, error) {
	problems, err := s.problemRepo.GetAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	detailProblems := make([]*entity.ProblemWithAnswerInformation, 0, len(problems))
	for _, problem := range problems {
		answers, err := s.answerRepo.FindByProblem(problem.ID, nil)
		if err != nil {
			return nil, err
		}

		unchecked := 0
		uncheckedNearOverdue := 0
		uncheckedOverdue := 0
		for _, answer := range answers {
			if answer.Point != nil {
				continue
			}

			unchecked += 1

			if now.After(answer.CreatedAt.Add(s.uncheckedNearOverdueThreshold)) {
				uncheckedNearOverdue += 1
			}

			if now.After(answer.CreatedAt.Add(s.uncheckedOverdueThreshold)) {
				uncheckedOverdue += 1
			}
		}

		detailProblems = append(detailProblems, &entity.ProblemWithAnswerInformation{
			Problem: *problem,

			Unchecked:            uint(unchecked),
			UncheckedNearOverdue: uint(unchecked),
			UncheckedOverdue:     uint(unchecked),
		})
	}

	return detailProblems, nil
}

func (s *ProblemService) FindByID(id uuid.UUID) (*entity.Problem, error) {
	return s.problemRepo.FindByID(id)
}

func (s *ProblemService) FindByCode(code string) (*entity.Problem, error) {
	return s.problemRepo.FindByCode(code)
}

func (s *ProblemService) Update(id uuid.UUID, req *UpdateProblemRequest) (*entity.Problem, error) {
	prob, err := s.problemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if prob == nil {
		return nil, errors.New("problem not found")
	}

	prob.AuthorID = req.AuthorID
	prob.Title = req.Title
	prob.Body = req.Body
	prob.Point = req.Point
	prob.PreviousProblemID = req.PreviousProblemID
	prob.SolvedCriterion = req.SolvedCriterion

	if err = prob.Validate(); err != nil {
		return nil, err
	}

	if req.PreviousProblemID != nil {
		prevProb, err := s.problemRepo.FindByID(*req.PreviousProblemID)
		if err != nil {
			return nil, err
		}
		if prevProb == nil {
			return nil, errors.New("previous_problem not found")
		}
	}

	return s.problemRepo.Update(prob)
}

func (s *ProblemService) Delete(id uuid.UUID) error {
	prob, err := s.problemRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.problemRepo.Delete(prob)
}
