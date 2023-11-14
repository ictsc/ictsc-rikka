package controller

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type ProblemController struct {
	problemService *service.ProblemService
}

func NewProblemController(problemService *service.ProblemService) *ProblemController {
	return &ProblemController{
		problemService: problemService,
	}
}

type CreateProblemRequest struct {
	Code              string     `json:"code"`
	Title             string     `json:"title"`
	Body              string     `json:"body"`
	Point             uint       `json:"point"`
	PreviousProblemID *uuid.UUID `json:"previous_problem_id"`
	SolvedCriterion   uint       `json:"solved_criterion"`
}

type CreateProblemResponse struct {
	Problem *entity.Problem `json:"problem"`
}

func (c *ProblemController) Create(req *CreateProblemRequest) (*CreateProblemResponse, error) {
	prob, err := c.problemService.Create(&service.CreateProblemRequest{
		Code:              req.Code,
		Title:             req.Title,
		Body:              req.Body,
		Point:             req.Point,
		PreviousProblemID: req.PreviousProblemID,
		SolvedCriterion:   req.SolvedCriterion,
	})

	if err != nil {
		return nil, err
	}
	return &CreateProblemResponse{
		Problem: prob,
	}, nil
}

type FindProblemByIDResponse struct {
	Problem *entity.Problem `json:"problem"`
}

func (c *ProblemController) FindByID(id string, metadataOnly bool) (*FindProblemByIDResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	prob, err := c.problemService.FindByID(uuid)
	if err != nil {
		return nil, err
	}

	if metadataOnly {
		prob.Body = ""
	}

	return &FindProblemByIDResponse{
		Problem: prob,
	}, nil
}

type FindProblemByCodeResponse struct {
	Problem *entity.Problem `json:"problem"`
}

func (c *ProblemController) FindByCode(code string, metadataOnly bool) (*FindProblemByCodeResponse, error) {
	prob, err := c.problemService.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if metadataOnly {
		prob.Body = ""
	}

	return &FindProblemByCodeResponse{
		Problem: prob,
	}, nil
}

type GetAllProblemsResponse struct {
	Problems []*entity.Problem `json:"problems"`
}

func (c *ProblemController) GetAll(metadataOnly bool) (*GetAllProblemsResponse, error) {
	probs, err := c.problemService.GetAll()
	if err != nil {
		return nil, err
	}

	if metadataOnly {
		for _, prob := range probs {
			prob.Body = ""
		}
	}

	return &GetAllProblemsResponse{
		Problems: probs,
	}, nil
}

type GetAllProblemsWithCurrentPointResponse struct {
	Problems []*entity.ProblemWithCurrentPoint `json:"problems"`
}

func (c *ProblemController) GetAllProblemsWithCurrentPoint(group *entity.UserGroup, metadataOnly bool) (*GetAllProblemsWithCurrentPointResponse, error) {
	probs, err := c.problemService.GetAllWithCurrentPoint(group)
	if err != nil {
		return nil, err
	}
	if metadataOnly {
		for _, prob := range probs {
			prob.Body = ""
		}
	}
	return &GetAllProblemsWithCurrentPointResponse{
		Problems: probs,
	}, nil
}

type GetAllProblemsWithAnswerInformationResponse struct {
	Problems []*entity.ProblemWithAnswerInformation `json:"problems"`
}

func (c *ProblemController) GetAllWithAnswerInformation(metadataOnly bool) (*GetAllProblemsWithAnswerInformationResponse, error) {
	probs, err := c.problemService.GetAllWithAnswerInformation()
	if err != nil {
		return nil, err
	}

	if metadataOnly {
		for _, prob := range probs {
			prob.Body = ""
		}
	}

	return &GetAllProblemsWithAnswerInformationResponse{
		Problems: probs,
	}, nil
}

type UpdateProblemRequest struct {
	Title             string     `json:"title"`
	Body              string     `json:"body"`
	Point             uint       `json:"point"`
	PreviousProblemID *uuid.UUID `json:"previous_problem_id"`
	SolvedCriterion   uint       `json:"solved_criterion"`
}

type UpdateProblemResponse struct {
	Problem *entity.Problem `json:"problem"`
}

func (c *ProblemController) Update(id string, req *UpdateProblemRequest) (*UpdateProblemResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	prob, err := c.problemService.Update(uuid, &service.UpdateProblemRequest{
		Title:             req.Title,
		Body:              req.Body,
		Point:             req.Point,
		PreviousProblemID: req.PreviousProblemID,
		SolvedCriterion:   req.SolvedCriterion,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateProblemResponse{
		Problem: prob,
	}, nil
}

func (c *ProblemController) Delete(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return c.problemService.Delete(uuid)
}
