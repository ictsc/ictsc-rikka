package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
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
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gopkg.in/yaml.v3"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
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
		&entity.Bastion{},
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
	if err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to open redis connection.").Error())
	}

	sameSiteMode := http.SameSiteDefaultMode
	switch config.Store.SameSiteStrictMode {
	case "lax":
		sameSiteMode = http.SameSiteLaxMode
	case "strict":
		sameSiteMode = http.SameSiteStrictMode
	case "none":
		sameSiteMode = http.SameSiteNoneMode
	}

	store.Options(sessions.Options{
		MaxAge:   43200,
		Path:     "/",
		Secure:   config.Store.Secure,
		HttpOnly: true,
		SameSite: sameSiteMode,
		Domain:   config.Store.Domain,
	})
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

	sqlDB.SetMaxOpenConns(config.MariaDB.MaxConn)
	sqlDB.SetMaxIdleConns(config.MariaDB.MaxConn)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}

func main() {
	userRepo := mariadb.NewUserRepository(db)
	userProfileRepo := mariadb.NewUserProfileRepository(db)
	userGroupRepo := mariadb.NewUserGroupRepository(db)
	bastionRepo := mariadb.NewBastionRepository(db)
	problemRepo := mariadb.NewProblemRepository(db)
	answerRepo := mariadb.NewAnswerRepository(db)
	attachmentRepo := mariadb.NewAttachmentRepository(db)
	s3Repo := s3repo.NewS3Repository(minioClient, config.Minio.BucketName)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, userProfileRepo, userGroupRepo)
	userGroupService := service.NewUserGroupService(userGroupRepo)
	bastionService := service.NewBastionService(bastionRepo)
	problemService := service.NewProblemService(config.Contest.AnswerLimit, userRepo, problemRepo, answerRepo)
	answerService := service.NewAnswerService(config.Contest.AnswerLimit, config.Notify.Answer, userRepo, answerRepo, problemRepo)
	attachmentService := service.NewAttachmentService(attachmentRepo, s3Repo)
	rankingService := service.NewRankingService(config.Contest.AnswerLimit, userGroupRepo, answerRepo)

	problemController := controller.NewProblemController(problemService)
	answerController := controller.NewAnswerController(answerService)
	attachmentController := controller.NewAttachmentController(attachmentService)

	errorMiddleware := middleware.NewErrorMiddleware()
	prometheus := ginprometheus.NewPrometheus("gin")
	if err := sentry.Init(sentry.ClientOptions{
		Debug:            true,
		Dsn:              config.Sentry.Dsn,
		Environment:      config.Sentry.Environment,
		TracesSampleRate: config.Sentry.TracesSampleRate,
	}); err != nil {
		log.Fatal(err.Error())
	}

	seed.Seed(&config.Seed, userRepo, userGroupRepo, *userService, *userGroupService, *bastionService, bastionRepo)

	r := gin.Default()

	origins := config.CORS.Origins
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = origins
	corsConfig.AllowCredentials = true

	r.Use(func(ctx *gin.Context) {
		rex := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
		uri := rex.ReplaceAllString(ctx.Request.RequestURI, "...")
		span := sentry.StartSpan(ctx.Request.Context(), "handleRequest",
			sentry.TransactionName(fmt.Sprintf("handleRequest: %s %s", ctx.Request.Method, uri)))
		ctx.Next()
		span.Finish()
	})

	r.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	r.Use(cors.New(corsConfig))
	r.Use(sessions.Sessions("session", store))

	prometheus.Use(r)

	api := r.Group("/api")
	api.Use(errorMiddleware.HandleError)
	{
		handler.NewAuthHandler(api, userRepo, authService, userService)
		handler.NewUserHandler(api, userRepo, userService)
		handler.NewUserGroupHandler(api, userRepo, userService, userGroupService, bastionService)
		handler.NewProblemHandler(api, userRepo, problemController, answerController)
		handler.NewAttachmentHandler(api, attachmentController, userRepo)
		handler.NewRankingHandler(api, userRepo, rankingService)

		api.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, "")
		})
	}

	go func() {
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
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	fmt.Println("signal received, quitting...")
	sqlDb, _ := db.DB()
	sqlDb.Close()
}
