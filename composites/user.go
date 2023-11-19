package composites

import (
	handlers "github.com/folklinoff/hack_and_change/controllers"
	"github.com/folklinoff/hack_and_change/repository"
	"gorm.io/gorm"
)

func NewUserComposite(db *gorm.DB) *handlers.UserHandler {
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	return userHandler
}
