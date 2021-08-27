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
	Body string `json:"body" binding:"required,max=10000"`
}

type CreateAnswerResponse struct {
	Answer *entity.Answer `json:"answer"`
}

func (c *AnswerController) Create(group *entity.UserGroup, problem_id string, req *CreateAnswerRequest) (*CreateAnswerResponse, error) {
	problem_uuid, err := uuid.Parse(problem_id)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.Create(&service.CreateAnswerRequest{
		UserGroup: group,
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

func (c *AnswerController) FindByID(group *entity.UserGroup, id string) (*FindByIDResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.FindByID(group, uuid)
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

// user group id is optional
func (c *AnswerController) FindByProblem(group *entity.UserGroup, probid string, userGroupId *uuid.UUID) (*FindByProblemResponse, error) {
	probuuid, err := uuid.Parse(probid)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.FindByProblem(group, probuuid, userGroupId)
	if err != nil {
		return nil, err
	}

	return &FindByProblemResponse{
		Answers: ans,
	}, nil
}

type FindByUserGroupResponse struct {
	Answers []*entity.Answer `json:"answers"`
}

type FindByUserGroupRequest struct {
	UserGroupID string `json:"user_group_id"`
}

func (c *AnswerController) FindByUserGroup(req *FindByUserGroupRequest) (*FindByUserGroupResponse, error) {
	uuid, err := uuid.Parse(req.UserGroupID)
	if err != nil {
		return nil, err
	}
	ans, err := c.answerService.FindByUserGroup(uuid)
	if err != nil {
		return nil, err
	}

	return &FindByUserGroupResponse{
		Answers: ans,
	}, nil
}

type FindByProblemAndUserGroupResponse struct {
	Answers []*entity.Answer `json:"answers"`
}

func (c *AnswerController) FindByProblemAndUserGroup(group *entity.UserGroup, probid string) (*FindByProblemAndUserGroupResponse, error) {
	probuuid, err := uuid.Parse(probid)
	if err != nil {
		return nil, err
	}

	ans, err := c.answerService.FindByProblemAndUserGroup(group, probuuid, group.ID)
	if err != nil {
		return nil, err
	}

	return &FindByProblemAndUserGroupResponse{
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
