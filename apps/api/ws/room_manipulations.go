package ws

import (
	"github.com/folklinoff/hack_and_change/models"
	"github.com/gorilla/websocket"
)

func CreateNewRoom(room models.Room) {
	rooms[room.ID] = room
	conns[room.ID] = make(map[int64]*websocket.Conn)
}

func AddUserToRoom(roomId int64, userId int64) {
	room := rooms[roomId]
	room.Users = append(room.Users, models.User{})
	// rooms[roomId].Users = append(rooms[roomId].Users)
	rooms[roomId] = room
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

func SendMessageToRoom(roomId int64, message models.Message) {
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
