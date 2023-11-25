package ws

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var rooms map[int64]models.Room
var conns map[int64]map[int64]*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		log.Println("CheckOrigin: ", r.Header.Get("Origin"))
		return true
	},
}

func ConnectToRoom(c *gin.Context) {
	c.Request.Header.Set("Upgrade", "Websocket")
	c.Request.Header.Set("Connection", "Upgrade")
	roomId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("ConnectToRoom: invalid room id")
		return
	}
	room := rooms[roomId]
	userId, err := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	if err != nil {
		log.Printf("ConnectToRoom: invalid user id")
		return
	}
	if !IsUserInTheRoom(roomId, userId) {
		log.Printf("ConnectToRoom: user %d is not in the room %d", userId, roomId)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ConnectToRoom: error occured when connecting to room: %s", err.Error())
		return
	}
	defer conn.Close()
	log.Printf("ConnectToRoom: new user joined with id: %d\n", userId)

	for i := range room.Messages {
		conn.WriteJSON(room.Messages[i])
	}
	log.Println("ConnectToRoom: all messages sent")
	ch := make(chan models.Message)
	go func() {
		for {
			var message models.Message
			conn.ReadJSON(message)
			ch <- message
			switch message.Type {
			case models.CreatedMessage:
				room.Messages = append(room.Messages, message)
			case models.JoinedMessage:
				room.Messages = append(room.Messages, message)
			case models.LeftMessage:

			}
		}
	}()
	for {
		select {
		case <-time.After(1 * time.Second):
			conn.WriteJSON(models.Message{Contents: "Hello, world!"})
		case msg := <-ch:
			conn.WriteJSON(models.Message{Contents: fmt.Sprintf("Received message %+v from client", msg)})
		}
		time.Sleep(1 * time.Second)
	}
}
