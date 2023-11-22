package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	if err := db.AutoMigrate(&repository.Room{}, &repository.User{}); err != nil {
		log.Fatalf("couldn't migrate into database: %s", err.Error())
	}

	if err := db.SetupJoinTable(&repository.Room{}, "UserRoom", &repository.User{}); err != nil {
		log.Fatalf("couldn't migrate into database: %s", err.Error())
	}

	e := gin.Default()
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "We gucci"})
	})
	e.GET("/room/:id", handlers.CORSMiddleware(), handlers.CheckAuthToken, ws.ConnectToRoom)

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
