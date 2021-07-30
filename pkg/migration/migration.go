package migration

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.UserGroup{},
		&entity.Problem{},
		&entity.Attachment{},
	)
}
