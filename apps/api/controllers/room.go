package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/folklinoff/hack_and_change/dto"
	"github.com/folklinoff/hack_and_change/models"
	"github.com/folklinoff/hack_and_change/ws"
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
	root.GET("/rooms", h.Get)
	root.POST("/rooms/:id", h.JoinRoom)
	root.GET("/rooms/:id", h.GetByID)
}

type RoomRepository interface {
	Save(ctx context.Context, room models.Room) (models.Room, error)
	GetByName(ctx context.Context, name string) (models.Room, error)
	GetByID(ctx context.Context, id int64) (models.Room, error)
	Get(ctx context.Context, name string) ([]models.Room, error)
	CheckPassword(ctx context.Context, roomId int64, password string) error
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
// @Param room body dto.CreateRoomRequestDTO true "Room body"
// @Param Set-Cookie header string true "Authorization token"
// @Success 201 {object} models.Room
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /rooms [post]
func (h *RoomHandler) Save(c *gin.Context) {
	ctx := c.Request.Context()
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	var req dto.CreateRoomRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind room: %s", err.Error())})
		return
	}
	room, err := h.repo.GetByName(ctx, req.Name)
	innerErr := errors.Unwrap(err)
	switch {
	case err == nil:
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: "room name is already taken"})
		return
	case innerErr != gorm.ErrRecordNotFound:
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get data about the room: %s", innerErr.Error())})
		return
	}
	room = dto.CreateRoomRequestDTOToModel(req)
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
	ws.CreateNewRoom(room)
	ws.AddUserToRoom(room.ID, user)
	c.JSON(http.StatusCreated, room)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Register for joining a room (not a websocket part)
// @Tags room
// @Accept  json
// @Produce  json
// @Param Set-Cookie header string true "Authorization token"
// @Param password body dto.JoinRoomRequestDTO false "Room with password"
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
		return
	}
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	repoRoom, err := h.repo.GetByID(ctx, roomId)
	if repoRoom.OwnerId == userId {
		c.JSON(http.StatusOK, repoRoom)
		return
	}
	for i := range repoRoom.Users {
		if repoRoom.Users[i].ID == userId {
			c.JSON(http.StatusOK, repoRoom)
			return
		}
	}
	if repoRoom.Private {
		var req dto.JoinRoomRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind room: %s", err.Error())})
			return
		}
		enteredPassword := req.Password
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get room by id: %s", err.Error())})
			return
		}
		err := h.repo.CheckPassword(ctx, roomId, enteredPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to log into room: %s", err.Error())})
			return
		}
	}
	user, err := h.userRepo.GetByID(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("failed to get user by id: %s", err.Error())})
		return
	}
	room, err := h.repo.GetByID(ctx, roomId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get room by id: %s", err.Error())})
		return
	}
	room.Users = append(room.Users, user)
	room, err = h.repo.Save(ctx, room)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to log into room: %s", err.Error())})
		return
	}
	ws.AddUserToRoom(roomId, user)
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
	roomId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("bad url: room id is not an integer: %s", err.Error())})
		return
	}
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)

	repoRoom, err := h.repo.GetByID(ctx, roomId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get room by id: %s", err.Error())})
		return
	}
	if repoRoom.OwnerId == userId {
		c.JSON(http.StatusOK, repoRoom)
		return
	}
	for i := range repoRoom.Users {
		if repoRoom.Users[i].ID == userId {
			c.JSON(http.StatusOK, repoRoom)
			return
		}
	}

	c.JSON(http.StatusBadRequest, ErrorMessage{Message: "you are not a member of this room"})
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
	name := c.Query("name")
	rooms, err := h.repo.Get(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get all rooms: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, rooms)
}
