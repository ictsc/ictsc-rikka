package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDBConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type db struct {
	config *MariaDBConfig
}

func newDB(config *MariaDBConfig) *db {
	return &db{
		config: config,
	}
}

func (d *db) init() (*gorm.DB, *sql.DB, error) {
	gormDB, err := gorm.Open(mysql.Open(d.getDSN()), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	conn, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}
	return gormDB, conn, nil
}

func (d *db) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.config.Username,
		d.config.Password,
		d.config.Address,
		d.config.Port,
		d.config.Database,
	)
}

func Migrate(config *MariaDBConfig) error {
	d := newDB(config)
	gormDB, conn, err := d.init()
	if err != nil {
		return err
	}
	defer conn.Close()

	return gormDB.AutoMigrate(
		&entity.User{},
		&entity.UserProfile{},
		&entity.UserGroup{},
		&entity.Problem{},
		&entity.Attachment{},
	)
}
