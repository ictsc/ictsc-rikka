package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ictsc/ictsc-rikka/pkg/repository/growi"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mariadb"
	"github.com/ictsc/ictsc-rikka/pkg/repository/rc"
	"github.com/ictsc/ictsc-rikka/pkg/service"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	configPath  string
	config      Config
	db          *gorm.DB
	redisClient *redis.Client
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

	redisClient, err = initRedis()
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "failed to init redis").Error())
	}
}

func initRedis() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Address, config.Redis.Port),
		Password: config.Redis.KeyPair,
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}

// TODO(k-shir0): Rikka とほぼ共通処理なので統一する
func initDatabase(c *MariaDBConfig) (*gorm.DB, error) {
	dsn := c.getDSN()

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}

func main() {
	ctx := context.Background()

	client := http.Client{}
	u, err := url.Parse(config.Growi.Url)
	if err != nil {
		panic(err)
	}

	growiSessionCookieRepo := rc.NewGrowiSessionCookieRepository(redisClient)
	problemWithSyncTimeRepo := rc.NewProblemWithSyncTimeRepository(redisClient)
	pageRepo := growi.NewPageRepository(&client, u, config.Growi.Token)
	subordinatedRepo := growi.NewSubordinatedPageRepository(&client, u, config.Growi.Path, config.Growi.Token)
	problemRepo := mariadb.NewProblemRepository(db)

	growiProblemSyncService := service.NewGrowiProblemSyncService(
		&client,
		u,
		config.Growi.Path,
		config.Growi.Username,
		config.Growi.Password,
		config.Growi.Token,
		growiSessionCookieRepo,
		problemWithSyncTimeRepo,
		pageRepo,
		subordinatedRepo,
		problemRepo,
	)

	err = growiProblemSyncService.Init(ctx)
	// TODO(k-shir0): エラー処理追加
	err = growiProblemSyncService.Sync(ctx)
	// TODO(k-shir0): エラー処理追加
	if err != nil {
		log.Fatal(err)
	}

}
