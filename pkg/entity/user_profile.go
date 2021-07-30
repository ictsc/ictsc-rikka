package entity

import "github.com/google/uuid"

type UserProfile struct {
	Base

	UserID           uuid.UUID `json:"user_id"`
	TwitterID        string    `json:"twitter_id"`
	GithubID         string    `json:"github_id"`
	FacebookID       string    `json:"facebook_id"`
	SelfIntroduction string    `json:"self_introduction"`
}
