package entity

import (
	"github.com/google/uuid"
)

type Attachment struct {
	Base Base
	User uuid.UUID
}
