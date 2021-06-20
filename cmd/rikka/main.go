package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/handler"
	"github.com/ictsc/ictsc-rikka/pkg/migration"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mariadb"
	"github.com/ictsc/ictsc-rikka/pkg/seed"
	"github.com/ictsc/ictsc-rikka/pkg/service"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"

	"github.com/gin-contrib/sessions"
	"gorm.io/gorm"
)

var (
	configPath string
	config     Config
	db         *gorm.DB
	store      redis.Store
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MariaDB.Username,
		config.MariaDB.Password,
		config.MariaDB.Address,
		config.MariaDB.Port,
		config.MariaDB.Database,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to open mariadb server.").Error())
	}

	if err := migration.Migrate(db); err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to migrate.").Error())
	}

	store, err = redis.NewStore(
		config.Redis.IdleConnectionSize,
		"tcp",
		fmt.Sprintf("%s:%d", config.Redis.Address, config.Redis.Port),
		config.Redis.Password,
		[]byte(config.Redis.KeyPair),
	)
	if err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to open redis connection.").Error())
	}

}

func main() {
	r := gin.Default()
	r.Use(sessions.Sessions("session", store))

	userRepo := mariadb.NewUserRepository(db)
	userGroupRepo := mariadb.NewUserGroupRepository(db)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, userGroupRepo)
	userGroupService := service.NewUserGroupService(userGroupRepo)

	seed.Seed(&config.Seed, userRepo, userGroupRepo, *userService, *userGroupService)

	api := r.Group("/api")
	{
		handler.NewAuthHandler(api, userRepo, authService, userService)
		handler.NewUserHandler(api, userRepo, userService)
		handler.NewUserGroupHandler(api, userRepo, userGroupService)
	}

	addr := fmt.Sprintf("%s:%d", config.Listen.Address, config.Listen.Port)
	if config.Listen.TLS == nil {
		// not serve tls
		if err := r.Run(addr); err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// serve tls
		if err := r.RunTLS(addr, config.Listen.TLS.CertFilePath, config.Listen.TLS.KeyFilePath); err != nil {
			log.Fatal(err.Error())
		}
	}

}
