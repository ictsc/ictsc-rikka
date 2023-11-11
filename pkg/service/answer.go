package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	e "github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type AnswerService struct {
	answerLimit time.Duration
	webhook     string
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

func NewAnswerService(answerLimit int, webhook string, userRepo repository.UserRepository, answerRepo repository.AnswerRepository, problemRepo repository.ProblemRepository) *AnswerService {
	return &AnswerService{
		answerLimit: time.Duration(answerLimit) * time.Minute,
		webhook:     webhook,
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
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("problem id is invalid")
	}

	if problem.Type == entity.MultipleType {
		type MultipleAnswer struct {
			Group int    `json:"group"`
			Value []int  `json:"value"`
			Type  string `json:"type"`
		}

		var answers []MultipleAnswer
		err := json.Unmarshal([]byte(ans.Body), &answers)
		if err != nil {
			return nil, e.NewBadRequestError("invalid answer format")
		}

		if len(answers) > len(problem.Answer) {
			return nil, e.NewBadRequestError("invalid answer format")
		}

		var sum uint
		for _, answer := range answers {
			group := answer.Group
			if group < 0 || group > len(problem.Answer) {
				return nil, e.NewBadRequestError("answer group is invalid")
			}

			question := problem.Answer[group]
			if question.Type == entity.RadioButton {
				if question.CorrectAnswers[0] == uint(answer.Value[0]) {
					sum += question.Scoring.Correct
				}
			}

			if question.Type == entity.CheckBox {
				correctCount := 0
				for _, val := range answer.Value {
					if contains(question.CorrectAnswers, uint(val)) {
						correctCount++
					}
				}

				if correctCount == len(question.CorrectAnswers) {
					sum += question.Scoring.Correct
				} else if question.Scoring.PartialCorrect != nil && correctCount > 0 {
					sum += *question.Scoring.PartialCorrect * uint(correctCount)
				}
			}
		}
		ans.Point = &sum
	}

	answer, err := s.answerRepo.Create(ans)
	if err != nil {
		return nil, err
	}

	//TODO: クリーンアーキテクチャ的にここでするべきではないので後でプレゼンターにする
	{
		text := fmt.Sprintf("<https://contest.mgmt.ictsc.net/scoring/%s?answer_id=%s |新着解答> 問題名:%s チーム名:%s",
			strings.ToLower(problem.Code), answer.ID, problem.Title, req.UserGroup.Name)

		param := struct {
			Text    string `json:"text"`
			Channel string `json:"channel"`
		}{
			Text:    text,
			Channel: "#problem-" + strings.ToLower(problem.Code),
		}
		json_str, err := json.Marshal(param)
		if err != nil {
			log.Println(err.Error())
			return answer, nil
		}
		resp, err := http.Post(s.webhook, "application/json", bytes.NewBuffer(json_str))
		if err != nil {
			log.Println(err.Error())
			return answer, nil
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	return answer, nil
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

func contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
