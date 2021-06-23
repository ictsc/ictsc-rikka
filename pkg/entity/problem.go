package entity

import "github.com/google/uuid"

type Problem struct {
	Base

	Code              string    `json:"code" gorm:"unique;index"`
	AuthorID          uuid.UUID `json:"author_id"`
	Author            *User     `json:"-"`
	Title             string    `json:"title"`
	Body              string    `json:"body,omitempty"`
	Point             uint      `json:"point"`
	PreviousProblemID *uuid.UUID `json:"previous_problem_id"`
	PreviousProblem   *Problem  `json:"-"`
	SolvedCriterion   uint      `json:"solved_criterion"`
}
