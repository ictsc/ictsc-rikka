package controller

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type AnswerController struct {
	answerService *service.AnswerService
}

func NewAnswerController(answerService *service.AnswerService) *AnswerController {
	return &AnswerController{
		answerService: answerService,
	}
}

type CreateAnswerRequest struct {
	Body string `json:"body"`
	Point uint `json:"point"`
	ProblemID *uuid.UUID `json:"problem_id"`
}

type CreateAnswerResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) Create(req *CreateAnswerRequest) (*CreateAnswerResponse, error) {
	ans, err := c.answerService.Create(&service.CreateAnswerRequest{
		Body: req.Body,
		Point: req.Point,
		ProblemID: req.ProblemID,
	})

	if err != nil {
		return nil, err
	}
	return &CreateAnswerResponse{
		Answer: ans,
	}, nil
}

type FindByIDResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) FindByID(id string) (*FindByIDResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.FindByID(uuid)
	if err != nil {
		return nil, err
	}

	return &FindByIDResponse{
		Answer: ans,
	}, nil
}

type FindByProblemResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) FindByProblem(id string) (*FindByProblemResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	ans, err := c.answerService.FindByProblem(uuid)
	if err != nil {
		return nil, err
	}

	return &FindByProblemResponse{
		Answer: ans,
	}, nil
}

type FindByTeamResponse struct {
	Answer *entity.Answer `json:"answer"`
}

type FindByTeamRequest struct {
	TeamID string `json:"team_group_id"`
}

func (c *AnswerController) FindByTeam(req *FindByTeamRequest) (*FindByTeamResponse, error) {
	uuid, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, err
	}
	ans, err := c.answerService.FindByTeam(uuid)
	if err != nil {
		return nil, err
	}

	return &FindByTeamResponse{
		Answer: ans,
	}, nil
}

type FindByProblemAndTeamResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) FindByProblemAndTeam(probid string, teamid string) (*FindByProblemAndTeamResponse, error) {
	probuuid, err := uuid.Parse(probid)
	if err != nil {
		return nil, err
	}
	teamuuid, err := uuid.Parse(teamid)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.FindByProblemAndTeam(probuuid,teamuuid)
	if err != nil {
		return nil, err
	}

	return &FindByProblemAndTeamResponse{
		Answer: ans,
	}, nil
}

type GetAllAnswersResponse struct {
	Answers []*entity.Answer `json:"answers"`
}

func (c *AnswerController) GetAll() (*GetAllAnswersResponse, error) {
	ans, err := c.answerService.GetAll()
	if err != nil {
		return nil, err
	}

	return &GetAllAnswersResponse{
		Answers: ans,
	}, nil
}

type UpdateAnswerRequest struct {
	Point             uint      `json:"point"`
	Body           string       `json:"body"`
	ProblemID *uuid.UUID `json:"problem_id"`
}

type UpdateAnswerResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) Update(id string, req *UpdateAnswerRequest) (*UpdateAnswerResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.Update(uuid, &service.UpdateAnswerRequest{
		Body: req.Body,
		Point: req.Point,
		ProblemID: req.ProblemID,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateAnswerResponse{
		Answer: ans,
	}, nil
}

func (c *AnswerController) Delete(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return c.answerService.Delete(uuid)
}
