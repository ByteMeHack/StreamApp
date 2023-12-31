package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/folklinoff/hack_and_change/composites"
	"github.com/folklinoff/hack_and_change/config"
	handlers "github.com/folklinoff/hack_and_change/controllers"
	"github.com/folklinoff/hack_and_change/database"
	_ "github.com/folklinoff/hack_and_change/docs"
	"github.com/folklinoff/hack_and_change/repository"
	"github.com/folklinoff/hack_and_change/ws"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/general_api_info.html
// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host bytemehack.ru
// @BasePath /api
func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("couldn't initialize database: %s", err.Error())
	}

	if err := db.AutoMigrate(&repository.Room{}, &repository.User{}, &repository.UserRoom{}); err != nil {
		log.Fatalf("couldn't migrate into database: %s", err.Error())
	}

	if err := db.SetupJoinTable(&repository.Room{}, "Users", &repository.UserRoom{}); err != nil {
		log.Fatalf("couldn't setup join table for room: %s", err.Error())
	}

	if err := db.SetupJoinTable(&repository.User{}, "Rooms", &repository.UserRoom{}); err != nil {
		log.Fatalf("couldn't setup join table for users: %s", err.Error())
	}
	ws.Init(db)

	e := gin.Default()
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// this endpoint is used to serve files to the client
	e.Static("/.well-known", "./static")
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "We gucci"})
	})
	e.GET("/room/:id", handlers.CORSMiddleware(), handlers.CheckAuthToken, ws.ConnectToRoom)
	// client := NewAwsClient()
	e.Use(handlers.CORSMiddleware())

	userHandler := composites.NewUserComposite(db)
	userHandler.Register(e)
	roomHandler := composites.NewRoomComposite(db)
	roomHandler.Register(e)

	go func() {
		log.Printf("server started on port :%s", config.Cfg.Port)
		err := http.ListenAndServe(fmt.Sprintf(":%s", config.Cfg.Port), e)
		if err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}

func NewAwsClient() *s3.Client {
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-central1" {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           "https://storage.yandexcloud.net",
				SigningRegion: "ru-central1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(awsCfg)
	return client
}
