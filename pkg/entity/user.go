package entity

import "github.com/google/uuid"

type User struct {
	Base

	Name           string
	DisplayName    string
	PasswordDigest string `json:"-"`
	UserGroupID    uuid.UUID
	UserGroup      *UserGroup `json:",omitempty"`
	IsReadOnly     bool
}
