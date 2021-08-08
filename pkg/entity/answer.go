package entity

import (
	"github.com/google/uuid"
)

type Answer struct {
	Base

	Point     uint      `json:"point"`
	Body      string    `json:"body" gorm:"not null"`
	Group     uuid.UUID `json:"group" gorm:"not null"`
	ProblemID uuid.UUID `json:"problem_id" gorm:"not null"`
}
