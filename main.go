package main

import (
	"fmt"
	"net/http"

	"github.com/folklinoff/hack_and_change/ws"
	"github.com/gin-gonic/gin"
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

// @host localhost:8080
// @BasePath /
func main() {
	e := gin.Default()
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "We gucci"})
	})
	e.GET("/room", ws.ConnectToRoom)
	fmt.Println("gucci")
	e.Run(":8080")
}
