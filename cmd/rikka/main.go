package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ictsc/ictsc-rikka/pkg/controller"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/handler"
	"github.com/ictsc/ictsc-rikka/pkg/delivery/http/middleware"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mariadb"
	"github.com/ictsc/ictsc-rikka/pkg/repository/s3repo"
	"github.com/ictsc/ictsc-rikka/pkg/seed"
	"github.com/ictsc/ictsc-rikka/pkg/service"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
)

var (
	configPath  string
	config      Config
	store       redis.Store
	minioClient *minio.Client
	db          *gorm.DB
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
		log.Fatalf(errors.Wrapf(err, "Failed to initialize database").Error())
	}

	db.AutoMigrate(
		&entity.User{},
		&entity.UserProfile{},
		&entity.UserGroup{},
		&entity.Problem{},
		&entity.Answer{},
		&entity.Attachment{},
	)

	store, err = redis.NewStore(
		config.Redis.IdleConnectionSize,
		"tcp",
		fmt.Sprintf("%s:%d", config.Redis.Address, config.Redis.Port),
		config.Redis.Password,
		[]byte(config.Redis.KeyPair),
	)
	store.Options(sessions.Options{
		MaxAge:   43200,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	if err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to open redis connection.").Error())
	}
	minioClient, err = minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKeyID, config.Minio.SecretAccessKey, ""),
		Secure: config.Minio.UseSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func initDatabase(c *MariaDBConfig) (*gorm.DB, error) {
	dsn := c.getDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(48)
	sqlDB.SetMaxIdleConns(48)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}

func main() {
	userRepo := mariadb.NewUserRepository(db)
	userProfileRepo := mariadb.NewUserProfileRepository(db)
	userGroupRepo := mariadb.NewUserGroupRepository(db)
	problemRepo := mariadb.NewProblemRepository(db)
	answerRepo := mariadb.NewAnswerRepository(db)
	attachmentRepo := mariadb.NewAttachmentRepository(db)
	s3Repo := s3repo.NewS3Repository(minioClient, config.Minio.BucketName)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, userProfileRepo, userGroupRepo)
	userGroupService := service.NewUserGroupService(userGroupRepo)
	problemService := service.NewProblemService(userRepo, problemRepo)
	answerService := service.NewAnswerService(userRepo, answerRepo, problemRepo)
	attachmentService := service.NewAttachmentService(attachmentRepo, s3Repo)
	rankingService := service.NewRankingService(userGroupRepo, answerRepo)

	problemController := controller.NewProblemController(problemService)
	answerController := controller.NewAnswerController(answerService)
	attachmentController := controller.NewAttachmentController(attachmentService)

	errorMiddleware := middleware.NewErrorMiddleware()

	seed.Seed(&config.Seed, userRepo, userGroupRepo, *userService, *userGroupService)

	r := gin.Default()

	origins := config.CORS.Origins
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = origins
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	r.Use(sessions.Sessions("session", store))

	api := r.Group("/api")
	api.Use(errorMiddleware.HandleError)
	{
		handler.NewAuthHandler(api, userRepo, authService, userService)
		handler.NewUserHandler(api, userRepo, userService)
		handler.NewUserGroupHandler(api, userRepo, userGroupService)
		handler.NewProblemHandler(api, userRepo, problemController, answerController)
		handler.NewAttachmentHandler(api, attachmentController, userRepo)
		handler.NewRankingHandler(api, userRepo, rankingService)
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
