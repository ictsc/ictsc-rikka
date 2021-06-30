package entity

import (
	"io"

	"github.com/google/uuid"
)

type Attachment struct {
	ID     uuid.UUID
	Name   string
	Reader io.Reader
}
