package main

import (
	"flag"
	"github.com/ictsc/growi_client"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mariadb"
	"github.com/ictsc/ictsc-rikka/pkg/service"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
)

var (
	configPath string
	config     Config
	db         *gorm.DB
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "config path")
	flag.Parse()

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "Failed to open config file `%s`.", configPath).Error())
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to decode config.").Error())
	}

	db, err = initDatabase(&config.MariaDB)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "failed to init database").Error())
	}
}

// TODO(k-shir0): Rikka とほぼ共通処理なので統一する
func initDatabase(c *MariaDBConfig) (*gorm.DB, error) {
	dsn := c.getDSN()

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}

func main() {
	u, err := url.Parse(config.Growi.Url)
	if err != nil {
		log.Fatal(err)
	}

	problemRepo := mariadb.NewProblemRepository(db)
	noticeRepo := mariadb.NewNoticeRepository(db)

	client := growi_client.NewGrowiClient(&growi_client.GrowiClientOption{
		URL:         u,
		AccessToken: config.Growi.Token,
	})
	growiProblemSyncService := service.NewGrowiProblemSyncService(
		client,
		config.Growi.ProblemPath,
		problemRepo,
	)
	growiNoticeSyncService := service.NewGrowiNoticeSyncService(
		client,
		config.Growi.NoticePath,
		noticeRepo,
	)

	if err := growiProblemSyncService.Sync(); err != nil {
		log.Fatal(err)
	}
	if err := growiNoticeSyncService.Sync(); err != nil {
		log.Fatal(err)
	}
}
