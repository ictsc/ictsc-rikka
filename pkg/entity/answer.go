package entity

import (
	"github.com/google/uuid"
)

type Answer struct {
	Base

	Point             uint      `json:"point"`
	Body           string       `json:"body"`
	ProblemID *uuid.UUID `json:"problem_id"`
}
