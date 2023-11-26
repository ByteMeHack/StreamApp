package ws

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/folklinoff/hack_and_change/repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var rooms map[int64]models.Room
var conns map[int64]map[int64]*websocket.Conn

var Repo RoomRepository

func Init(db *gorm.DB) {
	Repo = repository.NewRoomRepository(db)
	repoRooms, err := Repo.Get(context.Background(), "")
	if err != nil {
		log.Fatalf("init: couldn't get rooms: %s", err.Error())
		return
	}
	rooms = make(map[int64]models.Room)
	conns = make(map[int64]map[int64]*websocket.Conn)
	for i := range repoRooms {
		rooms[repoRooms[i].ID] = repoRooms[i]
		conns[repoRooms[i].ID] = make(map[int64]*websocket.Conn)
	}
}

type RoomRepository interface {
	Get(ctx context.Context, name string) ([]models.Room, error)
	Save(ctx context.Context, room models.Room) (models.Room, error)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
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
	room, ok := rooms[roomId]
	if !ok {
		log.Printf("ConnectToRoom: room with id %d doesn't exist", roomId)
		return
	}
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
		log.Printf("ConnectToRoom: error occured when connecasdfting to room: %s", err.Error())
		return
	}
	connDoneCh := make(chan bool)
	defer conn.Close()
	// In case user leaves the site without closing the connection
	go func() {
		<-c.Request.Context().Done()
		connDoneCh <- true
	}()

	// In case user closes the connection
	go func() {
		<-connDoneCh
		conn.Close()
		delete(conns[roomId], userId)
	}()

	log.Printf("ConnectToRoom: new user joined with id: %d\n", userId)
	conns[roomId][userId] = conn

	// Send all messages to the user
	for {
		conn.ReadMessage()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	for i := range room.Messages {
		conn.WriteJSON(room.Messages[i])
	}
	log.Println("ConnectToRoom: all messages sent")

	go ListenForIncomingMessages(connDoneCh, roomId, userId)

	// Test messages
	for {
		select {
		case <-connDoneCh:
			return
		default:
			conn.WriteJSON(
				models.Message{
					Contents:  "Hello, world!",
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				})
			time.Sleep(5 * time.Second)
		}
	}
}

func BroadcastMessageToRoom(roomId int64, message models.Message) {
	room := rooms[roomId]
	for i := range room.Users {
		log.Printf("BroadcastMessageToRoom: sending message to user %d", room.Users[i].ID)
		conn := conns[roomId][room.Users[i].ID]
		if conn == nil {
			continue
		}
		conn.WriteJSON(message)
	}
}

func ListenForIncomingMessages(connDoneCh chan bool, roomId int64, userId int64) {
	for {
		conn := conns[roomId][userId]
		if conn == nil {
			return
		}
		var message models.Message
		err := conn.ReadJSON(&message)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ListenForIncomingMessages: error occured when reading message: %s", err)
			}
		}()
		log.Printf("read message: %+v", message)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			connDoneCh <- true
			return
		}
		message.UserId = userId
		message.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		switch message.Type {
		case models.LeftMessage:
			DeleteUserFromRoom(roomId, userId)
			connDoneCh <- true
			Repo.Save(context.Background(), rooms[roomId])
		case models.KickMessage:
			if userId != rooms[roomId].OwnerId {
				log.Printf("ConnectToRoom: user %d is not the owner of the room %d", userId, roomId)
				continue
			}
			kickedUserId, err := strconv.ParseInt(message.Contents, 10, 64)
			if err != nil {
				log.Printf("ConnectToRoom: invalid user id")
				continue
			}

			conns[roomId][kickedUserId].Close()
			delete(conns[roomId], kickedUserId)
			DeleteUserFromRoom(roomId, kickedUserId)
			Repo.Save(context.Background(), rooms[roomId])
		case models.RegularMessage:

		}
		SaveMessageToRoom(roomId, message)
		BroadcastMessageToRoom(roomId, message)
	}
}
