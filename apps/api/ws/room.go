package ws

import (
	"log"
	"strconv"

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
	c.Request.ParseForm()
	roomId, err := strconv.ParseInt(c.Request.FormValue("id"), 10, 64)
	if err != nil {
		log.Printf("ConnectToRoom: invalid room id")
	}
	room := rooms[roomId]

	userId := c.Request.Header.Get("XUserID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ConnectToRoom: error occured when connecting to room: %s", err.Error())
		return
	}
	defer conn.Close()
	log.Printf("ConnectToRoom: new user joined with id: %s\n", userId)

	for i := range room.Messages {
		conn.WriteJSON(room.Messages[i])
	}
	log.Println("ConnectToRoom: all messages sent")

	for {
		var message models.Message
		conn.ReadJSON(message)

		switch message.Type {
		case models.CreatedMessage:
		}
	}
}

func CreateNewRoom(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Request.PostFormValue("id"), 10, 64)
	rooms[id] = models.Room{}
}
