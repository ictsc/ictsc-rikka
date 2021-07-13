package entity

import (
	"github.com/google/uuid"
)

type Answer struct {
	Base

	Point             uint      `json:"point"`
	Body           string       `json:"body"`
	Group             uuid.UUID      `json:"group"`
	ProblemID uuid.UUID `json:"problem_id"`
}
