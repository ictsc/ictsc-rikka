package entity

import (
	"github.com/google/uuid"
)

type Attachment struct {
	Base   Base
	UserID uuid.UUID
}
