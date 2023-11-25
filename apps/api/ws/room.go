package ws

import (
	"fmt"
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
	c.Request.ParseForm()
	roomId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("ConnectToRoom: invalid room id")
		return
	}
	room := rooms[roomId]

	userId := c.Request.Header.Get("XUserID")
	c.Request.Header.Set("Connection", "Upgrade")
	c.Request.Header.Set("Upgrade", "Websocket")
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
	ch := make(chan models.Message)
	go func() {
		for {
			var message models.Message
			conn.ReadJSON(message)
			ch <- message
			switch message.Type {
			case models.CreatedMessage:
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

func CreateNewRoom(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Request.PostFormValue("id"), 10, 64)
	rooms[id] = models.Room{}
}

// func AddUserToRoom(c *gin.Context) {
// 	c.Request.ParseForm()
// 	roomId, err := strconv.ParseInt(c.Request.FormValue("id"), 10, 64)
// 	if err != nil {
// 		log.Printf("AddUserToRoom: invalid room id")
// 	}
// 	room := rooms[roomId]

// 	userId := c.Request.Header.Get("XUserID")

// 	room.Users = append(room.Users, models.User{ID: userId})
// 	*rooms[roomId].Users = append(*rooms[roomId].Users)
// }
