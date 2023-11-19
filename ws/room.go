package ws

import (
	"log"
	"strconv"
	"time"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var rooms map[int64]models.Room

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectToRoom(c *gin.Context) {
	authToken, err := c.Cookie("XAuthorizationToken")
	if err != nil {
		log.Printf("ConnectToRoom: no token cookie provided: %s", err.Error())
		return
	}
	log.Println(authToken)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ConnectToRoom: error occured when connecting to room: %s", err.Error())
		return
	}
	defer conn.Close()
	log.Println("ConnectToRoom: new user joined")
	for {
		conn.WriteMessage(websocket.TextMessage, []byte("Hello, websocket!"))
		time.Sleep(time.Second)
		// conn.ReadJSON()
	}
}

func CreateNewRoom(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Request.PostFormValue("id"), 10, 64)
	rooms[id] = models.Room{}
}
