package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
)

type ProblemService struct {
	preRoundMode                  bool
	uncheckedNearOverdueThreshold time.Duration
	uncheckedOverdueThreshold     time.Duration

	userRepo    repository.UserRepository
	problemRepo repository.ProblemRepository
	answerRepo  repository.AnswerRepository
}

type CreateProblemRequest struct {
	Code              string
	Title             string
	Body              string
	Type              entity.ProblemType
	CorrectAnswers    []entity.CorrectAnswer
	Point             uint
	PreviousProblemID *uuid.UUID
	SolvedCriterion   uint
}

type UpdateProblemRequest struct {
	ID                uuid.UUID
	Code              string
	Title             string
	Body              string
	Type              entity.ProblemType
	Answer            []entity.CorrectAnswer
	Point             uint
	PreviousProblemID *uuid.UUID
	SolvedCriterion   uint
}

func NewProblemService(preRoundMode bool, answerLimit int, userRepo repository.UserRepository, problemRepo repository.ProblemRepository, answerRepo repository.AnswerRepository) *ProblemService {
	return &ProblemService{
		preRoundMode:                  preRoundMode,
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
		Title:             req.Title,
		Body:              req.Body,
		Type:              req.Type,
		CorrectAnswers:    req.CorrectAnswers,
		Point:             req.Point,
		PreviousProblemID: req.PreviousProblemID,
		SolvedCriterion:   req.SolvedCriterion,
	}

	if err := prob.DeleteMatterQuestionWithQuestionFieldAttach(); err != nil {
		return nil, err
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

// FullAccessの場合 groupIDはnullである
func (s *ProblemService) GetCurrentPointAnswer(problem *entity.Problem, groupId *uuid.UUID, IsfullAccess bool) (*entity.Answer, uint, uint, uint, error) {
	var currentAnswer *entity.Answer
	var answers []*entity.Answer
	var err error

	if IsfullAccess && groupId == nil {
		answers, err = s.answerRepo.FindByProblem(problem.ID, nil)
	} else {
		answers, err = s.answerRepo.FindByProblemAndUserGroup(problem.ID, *groupId)
	}
	if err != nil {
		return nil, 0, 0, 0, err
	}
	// まだユーザに点数が公開されていない回答を除外する
	now := time.Now()
	var CurrentPoint uint = 0
	var gotAt time.Time
	unchecked := 0
	uncheckedNearOverdue := 0
	uncheckedOverdue := 0
	for _, answer := range answers {

		// 20分制約によって回答が見れていない場合
		if !IsfullAccess && !now.After(answer.CreatedAt.Add(s.uncheckedOverdueThreshold)) {
			continue
		}

		if answer.Point == nil {
			unchecked += 1

			if now.After(answer.CreatedAt.Add(s.uncheckedNearOverdueThreshold)) {
				uncheckedNearOverdue += 1
			}

			if now.After(answer.CreatedAt.Add(s.uncheckedOverdueThreshold)) {
				uncheckedOverdue += 1
			}

			continue
		}

		if currentAnswer != nil {
			gotAt = currentAnswer.CreatedAt
		}

		//一番最新で一番高い得点かつ一番最新のanswerを出す
		if CurrentPoint < *answer.Point || (CurrentPoint == *answer.Point && answer.CreatedAt.Before(gotAt)) {
			currentAnswer = answer
			CurrentPoint = *answer.Point
			gotAt = answer.CreatedAt
		}
	}
	if currentAnswer == nil {
		return nil, uint(unchecked), uint(uncheckedNearOverdue), uint(uncheckedOverdue), errors.New("No answer yet")
	}
	return currentAnswer, uint(unchecked), uint(uncheckedNearOverdue), uint(uncheckedOverdue), nil
}

func (s *ProblemService) GetCurrentPoint(problem *entity.Problem, group *entity.UserGroup) uint {
	if answer, _, _, _, err := s.GetCurrentPointAnswer(problem, &group.ID, false); err == nil {
		return *answer.Point
	}
	return 0
}

func (s *ProblemService) GetCurrentPointWithUncheckedInformation(problem *entity.Problem) (*entity.ProblemWithAnswerInformation, error) {
	answer, unchecked, uncheckedNearOverdue, uncheckedOverdue, err := s.GetCurrentPointAnswer(problem, nil, true)
	if err != nil {
		// 全ての回答でNo answer yetの場合か、その他エラーの場合
		return &entity.ProblemWithAnswerInformation{
			Problem: *problem,

			Unchecked:            unchecked,
			UncheckedNearOverdue: uncheckedNearOverdue,
			UncheckedOverdue:     uncheckedOverdue,
			CurrentPoint:         0,
			IsSolved:             false,
		}, nil
	}
	return &entity.ProblemWithAnswerInformation{
		Problem: *problem,

		Unchecked:            unchecked,
		UncheckedNearOverdue: uncheckedNearOverdue,
		UncheckedOverdue:     uncheckedOverdue,
		CurrentPoint:         *answer.Point,
		IsSolved:             *answer.Point >= problem.SolvedCriterion,
	}, nil

}

func (s *ProblemService) GetAll() ([]*entity.Problem, error) {
	return s.problemRepo.GetAll()
}

func (s *ProblemService) GetAllWithCurrentPoint(group *entity.UserGroup) ([]*entity.ProblemWithCurrentPoint, error) {
	problems, err := s.problemRepo.GetProblemsWithIsAnsweredByUserGroup(group.ID)
	if err != nil {
		return nil, err
	}
	detailProblems := make([]*entity.ProblemWithCurrentPoint, 0, len(problems))
	for _, problem := range problems {
		var CurrentPoint uint
		if !s.preRoundMode {
			CurrentPoint = s.GetCurrentPoint(&problem.Problem, group)
		}
		detailProblems = append(detailProblems, &entity.ProblemWithCurrentPoint{
			Problem: problem.Problem,

			IsAnswered:   problem.IsAnswered,
			CurrentPoint: CurrentPoint,
			IsSolved:     CurrentPoint >= problem.SolvedCriterion,
		})
	}
	return detailProblems, nil
}

func (s *ProblemService) GetAllWithAnswerInformation() ([]*entity.ProblemWithAnswerInformation, error) {
	problems, err := s.problemRepo.GetAll()
	if err != nil {
		return nil, err
	}

	detailProblems := make([]*entity.ProblemWithAnswerInformation, 0, len(problems))
	for _, problem := range problems {
		problemWithAnswerInfo, _ := s.GetCurrentPointWithUncheckedInformation(problem)
		detailProblems = append(detailProblems, problemWithAnswerInfo)
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

	return s.problemRepo.Update(prob, false)
}

func (s *ProblemService) Delete(id uuid.UUID) error {
	prob, err := s.problemRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.problemRepo.Delete(prob)
}
