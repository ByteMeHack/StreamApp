package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomHandler struct {
	repo     RoomRepository
	userRepo UserRepository
}

func (h *RoomHandler) Register(e *gin.Engine) {
	root := e.Group("/", CORSMiddleware(), CheckAuthToken)
	root.POST("/rooms", h.Save)
	root.OPTIONS("/rooms")
	root.GET("/rooms", h.Get)
	root.OPTIONS("rooms/:id")
	root.POST("/rooms/:id", h.JoinRoom)
	root.GET("/rooms/:id", h.GetByID)
}

type RoomRepository interface {
	Save(ctx context.Context, room models.Room) (models.Room, error)
	GetByName(ctx context.Context, name string) (models.Room, error)
	GetByID(ctx context.Context, id int64) (models.Room, error)
	Get(ctx context.Context) ([]models.Room, error)
	LogIntoRoom(ctx context.Context, id, userId int64, password string) (models.Room, error)
}

func NewRoomHandler(repo RoomRepository, userRepo UserRepository) *RoomHandler {
	return &RoomHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Create new room
// @Tags room
// @Accept  json
// @Produce  json
// @Param room body models.Room true "Room body"
// @Param Set-Cookie header string true "Authorization token"
// @Success 201 {object} models.Room
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /rooms [post]
func (h *RoomHandler) Save(c *gin.Context) {
	ctx := c.Request.Context()
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind room: %s", err.Error())})
		return
	}
	_, err := h.repo.GetByName(ctx, room.Name)
	innerErr := errors.Unwrap(err)
	switch {
	case err == nil:
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: "room name is already taken"})
		return
	case innerErr != gorm.ErrRecordNotFound:
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get data about the room: %s", innerErr.Error())})
		return
	}
	room.OwnerId = userId
	user, err := h.userRepo.GetByID(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("failed to get user by id: %s", err.Error())})
		return
	}
	room.Users = append(room.Users, user)
	room, err = h.repo.Save(ctx, room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("failed to create room: %s", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, room)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Register for joining a room (not a websocket part)
// @Tags room
// @Accept  json
// @Produce  json
// @Param Set-Cookie header string true "Authorization token"
// @Success 200 {object} models.Room
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /rooms/:id [POST]
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	ctx := c.Request.Context()
	roomId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("bad url: room id is not an integer: %s", err.Error())})
	}
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	repoRoom, err := h.repo.GetByID(ctx, roomId)
	var room models.Room
	if repoRoom.Private {
		if err := c.ShouldBindJSON(&room); err != nil {
			c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind room: %s", err.Error())})
			return
		}
	}
	enteredPassword := room.Password
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get room by id: %s", err.Error())})
		return
	}
	if room.OwnerId == userId {
		c.JSON(http.StatusOK, room)
		return
	}
	for i := range room.Users {
		if room.Users[i].ID == userId {
			c.JSON(http.StatusOK, room)
			return
		}
	}
	room, err = h.repo.LogIntoRoom(ctx, roomId, userId, enteredPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to log into room: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, room)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get room by id
// @Tags room
// @Accept  json
// @Produce  json
// @Param Set-Cookie header string true "Authorization token"
// @Success 200 {object} models.Room
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /rooms/:id [get]
func (h *RoomHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	c.Request.ParseForm()
	roomId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("bad url: room id is not an integer: %s", err.Error())})
	}
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind room: %s", err.Error())})
		return
	}
	room, err = h.repo.GetByID(ctx, roomId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get room by id: %s", err.Error())})
		return
	}
	if room.OwnerId == userId {
		c.JSON(http.StatusOK, room)
		return
	}
	for i := range room.Users {
		if room.Users[i].ID == userId {
			c.JSON(http.StatusOK, room)
			return
		}
	}
	if room.Private {
		room, err = h.repo.LogIntoRoom(ctx, roomId, userId, room.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to log into room: %s", err.Error())})
			return
		}
	}
	c.JSON(http.StatusOK, room)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get all rooms
// @Tags room
// @Accept  json
// @Produce  json
// @Param Set-Cookie header string true "Authorization token"
// @Success 200 {object} []models.Room
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /rooms [get]
func (h *RoomHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	rooms, err := h.repo.Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get all rooms: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, rooms)
}
