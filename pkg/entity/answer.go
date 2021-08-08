package entity

import (
	"github.com/google/uuid"
)

type Answer struct {
	Base

	Point       *uint     `json:"point"`
	Body        string    `json:"body" gorm:"not null"`
	UserGroupID uuid.UUID `json:"user_group_id" gorm:"not null"`
	UserGroup   UserGroup `json:"user_group,omitempty"`
	ProblemID   uuid.UUID `json:"problem_id" gorm:"not null"`
	Problem     Problem   `json:"problem,omitempty"`
}
