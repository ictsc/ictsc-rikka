package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"size:191"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}
