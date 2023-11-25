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
	defer conn.Close()

	log.Printf("ConnectToRoom: new user joined with id: %d\n", userId)
	conns[roomId][userId] = conn
	for i := range room.Messages {
		conn.WriteJSON(room.Messages[i])
	}
	log.Println("ConnectToRoom: all messages sent")

	go ListenForIncomingMessages(roomId, userId, conn)
	for {
		<-time.After(1 * time.Second)
		conn.WriteJSON(models.Message{Contents: "Hello, world!"})
		time.Sleep(1 * time.Second)
	}
}

func BroadcastMessageToRoom(roomId int64, message models.Message) {
	room := rooms[roomId]
	for i := range room.Users {
		conn := conns[roomId][room.Users[i].ID]
		if conn == nil {
			log.Printf("BroadcastMessageToRoom: user with id %d is not connected to room %d", room.Users[i].ID, roomId)
			continue
		}
		conn.WriteJSON(message)
	}
}

func ListenForIncomingMessages(roomId int64, userId int64, conn *websocket.Conn) {
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Printf("ConnectToRoom: connection closed by user %d", userId)
			delete(conns[roomId], userId)
			return
		}
		message.UserId = userId
		switch message.Type {
		case models.LeftMessage:
			DeleteUserFromRoom(roomId, userId)
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
			DeleteUserFromRoom(roomId, kickedUserId)
			Repo.Save(context.Background(), rooms[roomId])
		}
		SaveMessageToRoom(roomId, message)
		BroadcastMessageToRoom(roomId, message)
	}
}
