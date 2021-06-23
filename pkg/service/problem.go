package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type ProblemService struct {
	userRepo      repository.UserRepository
	problemRepo repository.ProblemRepository
}

type CreateProblemRequest struct {
	Code string
	AuthorID uuid.UUID
	Title string
	Body string
	Point uint
	PreviousProblemID *uuid.UUID
	SolvedCriterion uint
}

type UpdateProblemRequest struct {
	ID   uuid.UUID
	Code string
	AuthorID uuid.UUID
	Title string
	Body string
	Point uint
	PreviousProblemID *uuid.UUID
	SolvedCriterion uint
}

func NewProblemService(userRepo repository.UserRepository, problemRepo repository.ProblemRepository) *ProblemService {
	return &ProblemService{
		userRepo:      userRepo,
		problemRepo: problemRepo,
	}
}

func (s *ProblemService) Create(req CreateProblemRequest) (*entity.Problem, error) {
	if req.PreviousProblemID != nil {
		prob, err := s.problemRepo.FindByID(*req.PreviousProblemID)
		if err != nil {
			return nil, err
		}
		if prob == nil {
			return nil, errors.New("previous_problem not found")
		}
	}

	// TODO: validate problem code

	if req.Point < req.SolvedCriterion {
		return nil, errors.New("solved_criterion must less than or equal to point")
	}

	return s.problemRepo.Create(&entity.Problem{
		Code: req.Code,
		AuthorID: req.AuthorID,
		Title: req.Title,
		Body: req.Body,
		Point: req.Point,
		PreviousProblemID: req.PreviousProblemID,
		SolvedCriterion: req.SolvedCriterion,
	})
}

func (s *ProblemService) GetAll() ([]*entity.Problem, error) {
	return s.problemRepo.GetAll()
}

func (s *ProblemService) FindByID(id uuid.UUID) (*entity.Problem, error) {
	return s.problemRepo.FindByID(id)
}

func (s *ProblemService) FindByCode(code string) (*entity.Problem, error) {
	return s.problemRepo.FindByCode(code)
}

func (s *ProblemService) Update(req UpdateProblemRequest) (*entity.Problem, error) {
	if req.PreviousProblemID != nil {
		prob, err := s.problemRepo.FindByID(*req.PreviousProblemID)
		if err != nil {
			return nil, err
		}
		if prob == nil {
			return nil, errors.New("previous_problem not found")
		}
	}

	// TODO: validate problem code

	return s.problemRepo.Update(&entity.Problem{
		Code: req.Code,
		AuthorID: req.AuthorID,
		Title: req.Title,
		Body: req.Body,
		Point: req.Point,
		PreviousProblemID: req.PreviousProblemID,
		SolvedCriterion: req.SolvedCriterion,
	})
}

func (s *ProblemService) Delete(id uuid.UUID) error {
	prob, err := s.problemRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.problemRepo.Delete(prob)
}
