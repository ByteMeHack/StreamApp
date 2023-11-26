package ws

import (
	"fmt"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/gorilla/websocket"
)

func CreateNewRoom(room models.Room) {
	rooms[room.ID] = room
	conns[room.ID] = make(map[int64]*websocket.Conn)
}

func AddUserToRoom(roomId int64, user models.User) {
	room := rooms[roomId]
	room.Users = append(room.Users, user)
	rooms[roomId] = room
	BroadcastMessageToRoom(
		roomId,
		models.Message{
			UserId:   user.ID,
			Type:     models.JoinedMessage,
			Contents: fmt.Sprintf("User with id %d joined the room", user.ID),
		})
}

func IsUserInTheRoom(roomId int64, userId int64) bool {
	room := rooms[roomId]
	for i := range room.Users {
		if room.Users[i].ID == userId {
			return true
		}
	}
	return false
}

func SaveMessageToRoom(roomId int64, message models.Message) {
	room := rooms[roomId]
	room.Messages = append(room.Messages, message)
	rooms[roomId] = room
}

func DeleteUserFromRoom(roomId int64, userId int64) {
	room := rooms[roomId]
	for i := range room.Users {
		if room.Users[i].ID == userId {
			room.Users = append(room.Users[:i], room.Users[i+1:]...)
			rooms[roomId] = room
			return
		}
	}
}
