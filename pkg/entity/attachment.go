package entity

import (
	"io"

	"github.com/google/uuid"
)

type Attachment struct {
	ID          uuid.UUID
	Reader      io.Reader
	User        string
	DisplayName string
	Password    string
	Group       string
	ReadOnly    bool
}
