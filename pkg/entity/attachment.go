package entity

import (
	"github.com/google/uuid"
)

type Attachment struct {
	Base
	UserID uuid.UUID
}
