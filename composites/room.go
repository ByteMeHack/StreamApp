package composites

import (
	handlers "github.com/folklinoff/hack_and_change/controllers"
	"github.com/folklinoff/hack_and_change/repository"
	"gorm.io/gorm"
)

func NewRoomComposite(db *gorm.DB) *handlers.RoomHandler {
	roomRepo := repository.NewRoomRepository(db)
	roomHandler := handlers.NewRoomHandler(roomRepo)
	return roomHandler
}
