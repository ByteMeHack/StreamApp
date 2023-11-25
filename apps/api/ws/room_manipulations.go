package ws

import (
	"log"
	"strconv"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/gin-gonic/gin"
)

func CreateNewRoom(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Request.PostFormValue("id"), 10, 64)
	rooms[id] = models.Room{}
}

func AddUserToRoom(c *gin.Context) {
	c.Request.ParseForm()
	roomId, err := strconv.ParseInt(c.Request.FormValue("id"), 10, 64)
	if err != nil {
		log.Printf("AddUserToRoom: invalid room id")
	}
	room := rooms[roomId]

	// userId := c.Request.Header.Get("XUserID")

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
