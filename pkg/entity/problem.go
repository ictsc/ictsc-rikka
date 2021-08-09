package entity

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"regexp"
)

type Problem struct {
	Base

	Code              string     `json:"code" gorm:"unique;index"`
	AuthorID          uuid.UUID  `json:"author_id"`
	Author            *User      `json:"-"`
	Title             string     `json:"title"`
	Body              string     `json:"body,omitempty"`
	Point             uint       `json:"point"`
	PreviousProblemID *uuid.UUID `json:"previous_problem_id"`
	PreviousProblem   *Problem   `json:"-"`
	SolvedCriterion   uint       `json:"solved_criterion"`
}

func (p *Problem) Validate() error {
	if matches, err := regexp.Match("[A-Z]{3}", []byte(p.Code)); err != nil {
		return err
	} else if !matches {
		return errors.New("code must match the pattern [A-Z]{3}")
	}

	if !(p.SolvedCriterion <= p.Point) {
		return errors.New("solved_criterion must be less than or equal to point")
	}

	return nil
}
