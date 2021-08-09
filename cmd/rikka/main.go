package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ictsc/ictsc-rikka/pkg/controller"

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

	if err := mariadb.Migrate(&config.MariaDB); err != nil {
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
	store.Options(sessions.Options{
		MaxAge:   3600,
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

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.CORS.Origins
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	r.Use(sessions.Sessions("session", store))

	errorMiddleware := middleware.NewErrorMiddleware()

	userRepo := mariadb.NewUserRepository(&config.MariaDB)
	userProfileRepo := mariadb.NewUserProfileRepository(&config.MariaDB)
	userGroupRepo := mariadb.NewUserGroupRepository(&config.MariaDB)
	problemRepo := mariadb.NewProblemRepository(&config.MariaDB)
	answerRepo := mariadb.NewAnswerRepository(&config.MariaDB)
	attachmentRepo := mariadb.NewAttachmentRepository(&config.MariaDB)
	s3Repo := s3repo.NewS3Repository(minioClient, config.Minio.BucketName)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, userProfileRepo, userGroupRepo)
	userGroupService := service.NewUserGroupService(userGroupRepo)
	problemService := service.NewProblemService(userRepo, problemRepo)
	answerService := service.NewAnswerService(userRepo, answerRepo, problemRepo)
	attachmentService := service.NewAttachmentService(attachmentRepo, s3Repo)

	problemController := controller.NewProblemController(problemService)
	answerController := controller.NewAnswerController(answerService)
	attachmentController := controller.NewAttachmentController(attachmentService)

	seed.Seed(&config.Seed, userRepo, userGroupRepo, *userService, *userGroupService)

	api := r.Group("/api")
	api.Use(errorMiddleware.HandleError)
	{
		handler.NewAuthHandler(api, userRepo, authService, userService)
		handler.NewUserHandler(api, userRepo, userService)
		handler.NewUserGroupHandler(api, userRepo, userGroupService)
		handler.NewProblemHandler(api, userRepo, problemController, answerController)
		handler.NewAttachmentHandler(api, attachmentController, userRepo)
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
