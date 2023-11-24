package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	auth "github.com/folklinoff/hack_and_change/pkg/auth"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CheckAuthToken(c *gin.Context) {
	uniqueKey := c.Request.Header.Get("App-Key")
	if uniqueKey != "" {
		if uniqueKey != os.Getenv("APP_KEY") {
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		} else {
			c.Next()
			return
		}
	}
	token, err := c.Cookie("XAuthorizationToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorMessage{Message: "No token cookie provided"})
		c.Abort()
		return
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, ErrorMessage{Message: "No token provided"})
		c.Abort()
		return
	}
	id, err := auth.Auth(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorMessage{Message: fmt.Sprintf("Couldn't authorize you: %s", err)})
		c.Abort()
		return
	}
	log.Println(token)
	log.Println("User id: ", id)
	c.Request.Header.Add("XUserID", strconv.FormatInt(id, 10))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := c.Request.Header["Origin"]
		origin := ""
		if len(origins) > 0 {
			origin = origins[0]
		}
		fmt.Println(origin)
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, App-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
