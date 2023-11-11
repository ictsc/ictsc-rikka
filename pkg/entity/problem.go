package entity

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Problem struct {
	Base

	Code              string          `json:"code" gorm:"unique;index"`
	AuthorID          uuid.UUID       `json:"author_id"`
	Author            *User           `json:"-"`
	Title             string          `json:"title"`
	Body              string          `json:"body,omitempty"`
	Type              ProblemType     `yaml:"type"`
	CorrectAnswers    []CorrectAnswer `json:"correct_answers,omitempty"`
	Point             uint            `json:"point"`
	PreviousProblemID *uuid.UUID      `json:"previous_problem_id"`
	PreviousProblem   *Problem        `json:"-"`
	SolvedCriterion   uint            `json:"solved_criterion"`
}

func (p *Problem) Validate() error {
	if matches, err := regexp.Match("[a-zA-Z]{3}", []byte(p.Code)); err != nil {
		return err
	} else if !matches {
		return errors.New("code must match the pattern [A-Z]{3}")
	}

	if !(p.SolvedCriterion <= p.Point) {
		return errors.New("solved_criterion must be less than or equal to point")
	}

	return nil
}

type ProblemWithAnswerInformation struct {
	Problem

	Unchecked            uint `json:"unchecked"`
	UncheckedNearOverdue uint `json:"unchecked_near_overdue"`
	UncheckedOverdue     uint `json:"unchecked_overdue"`

	CurrentPoint uint `json:"current_point"`
	IsSolved     bool `json:"is_solved"`
}

type ProblemWithCurrentPoint struct {
	Problem

	CurrentPoint uint `json:"current_point"`
	IsSolved     bool `json:"is_solved"`
}

type ProblemWithSyncTime struct {
	Problem
	UpdatedAt time.Time
}
