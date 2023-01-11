package entity

import "time"

type ProblemWithInfo struct {
	Problem   Problem
	UpdatedAt time.Time
}
