package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"size:191" json:"id"`

	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}
