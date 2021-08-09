package entity

import "github.com/google/uuid"

type User struct {
	Base

	Name           string       `json:"name" gorm:"unique;not null"`
	DisplayName    string       `json:"display_name" gorm:"not null"`
	PasswordDigest string       `json:"-" gorm:"not null"`
	UserGroupID    uuid.UUID    `json:"user_group_id"`
	UserGroup      *UserGroup   `json:"user_group,omitempty"`
	UserProfile    *UserProfile `json:"user_profile,omitempty"`
	IsReadOnly     bool         `json:"is_read_only" gorm:"not null"`
}
