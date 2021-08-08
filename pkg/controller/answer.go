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
}

type CreateAnswerResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) Create(problem_id string, groupuuid uuid.UUID, req *CreateAnswerRequest) (*CreateAnswerResponse, error) {
	problem_uuid, err := uuid.Parse(problem_id)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.Create(&service.CreateAnswerRequest{
		Group:     groupuuid,
		Body:      req.Body,
		ProblemID: problem_uuid,
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
	Answers []*entity.Answer `json:"answers"`
}

// team id is optional
func (c *AnswerController) FindByProblem(probid string, teamid string) (*FindByProblemResponse, error) {
	probuuid, err := uuid.Parse(probid)
	if err != nil {
		return nil, err
	}

	var teamuuid *uuid.UUID

	if teamid != "" {
		id, err := uuid.Parse(teamid)
		if err != nil {
			return nil, err
		}
		teamuuid = &id
	}

	ans, err := c.answerService.FindByProblem(probuuid, teamuuid)
	if err != nil {
		return nil, err
	}

	return &FindByProblemResponse{
		Answers: ans,
	}, nil
}

type FindByTeamResponse struct {
	Answers []*entity.Answer `json:"answers"`
}

type FindByTeamRequest struct {
	TeamID string `json:"team_group_id"`
}

func (c *AnswerController) FindByUserGroup(req *FindByTeamRequest) (*FindByTeamResponse, error) {
	uuid, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, err
	}
	ans, err := c.answerService.FindByUserGroup(uuid)
	if err != nil {
		return nil, err
	}

	return &FindByTeamResponse{
		Answers: ans,
	}, nil
}

type FindByProblemAndTeamResponse struct {
	Answers []*entity.Answer `json:"answers"`
}

func (c *AnswerController) FindByProblemAndUserGroup(probid string, teamuuid uuid.UUID) (*FindByProblemAndTeamResponse, error) {
	probuuid, err := uuid.Parse(probid)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.FindByProblemAndUserGroup(probuuid, teamuuid)
	if err != nil {
		return nil, err
	}

	return &FindByProblemAndTeamResponse{
		Answers: ans,
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
	Point uint `json:"point"`
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
		Point: req.Point,
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
