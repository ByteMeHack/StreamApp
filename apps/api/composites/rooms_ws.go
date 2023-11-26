package composites

import (
	"github.com/folklinoff/hack_and_change/repository"
	"github.com/folklinoff/hack_and_change/ws"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	repo := repository.NewRoomRepository(db)
	ws.Repo = repo
}
