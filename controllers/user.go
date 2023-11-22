package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/folklinoff/hack_and_change/dto"
	"github.com/folklinoff/hack_and_change/models"
	auth "github.com/folklinoff/hack_and_change/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (u *UserHandler) Register(e *gin.Engine) {
	root := e.Group("/", CORSMiddleware())
	root.OPTIONS("signup")
	root.POST("signup", u.SignUp)
	root.OPTIONS("login")
	root.POST("login", u.Login)
	root.OPTIONS("logout")
	root.POST("logout", u.Logout)
	root.OPTIONS("me")
	root.GET("me", CheckAuthToken, u.Me)
	root.OPTIONS("users/:id")
	root.GET("users/:id", CheckAuthToken, u.GetByID)
	root.OPTIONS("users/:id/rooms")
	root.GET("users/:id/rooms", CheckAuthToken, u.GetUserRooms)
}

type UserRepository interface {
	Save(ctx context.Context, user models.User) (models.User, error)
	GetByMail(ctx context.Context, mail string) (models.User, error)
	GetByID(ctx context.Context, id int64) (models.User, error)
	LoginUserByMail(ctx context.Context, mail string, password string) (models.User, error)
	GetUserRooms(ctx context.Context, id int64) ([]models.Room, error)
}

type LoginResponseDTO struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Create an account
// @Description Sign up using email, password
// @Tags accounts
// @Accept  json
// @Param user body dto.SignupRequestDTO true "User credentials"
// @Produce  json
// @Success 201 {object} dto.UserModelResponseDTO
// @Failure 400 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.SignupRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind user model: %s", err.Error())})
		return
	}
	user := dto.SignupRequestDTOToUserModel(req)
	_, err := h.repo.GetByMail(ctx, user.Email)
	innerErr := errors.Unwrap(err)
	switch {
	case err == nil:
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: "email is already taken"})
		return
	case innerErr != gorm.ErrRecordNotFound:
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get data about user: %s", innerErr.Error())})
		return
	}
	user, err = h.repo.Save(ctx, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("couldn't create user: %s", err.Error())})
		return
	}
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("couldn't create token: %s", err.Error())})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "XAuthorizationToken",
		Value:  token,
		Domain: strings.Split(c.Request.Host, ":")[0],
		Path:   "/",

		MaxAge: 36000,
	})
	c.JSON(http.StatusCreated, user)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Login into an account
// @Description Log in using email and password
// @Tags accounts
// @Accept  json
// @Param user body dto.LoginRequestDTO true "User credentials"
// @Produce  json
// @Success 200 {object} dto.UserModelResponseDTO
// @Failure 400 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind user model: %s", err.Error())})
		return
	}
	user := dto.LoginRequestDTOToUserModel(req)
	log.Println("User attempted to log in: ", user)
	_, err := h.repo.GetByMail(ctx, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("couldn't find user by mail: %s", err.Error())})
		return
	}
	user, err = h.repo.LoginUserByMail(ctx, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("incorrect password: %s", err.Error())})
		return
	}
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't create token: %s", err.Error())})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "XAuthorizationToken",
		Value:  token,
		Domain: strings.Split(c.Request.Host, ":")[0],
		Path:   "/",
		MaxAge: 36000,
	})
	c.JSON(http.StatusOK, user)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Logout from an account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Success 200 {object} object{message=string}
// @Router /accounts/{id} [get]
func (h *UserHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "XAuthorizationToken",
		Value:  "",
		Domain: strings.Split(c.Request.Host, ":")[0],
		Path:   "/",
		MaxAge: -1,
	})
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get general account info
// @Tags accounts
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.UserModelResponseDTO
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /me [get]
func (h *UserHandler) Me(c *gin.Context) {
	ctx := c.Request.Context()
	userId, _ := strconv.ParseInt(c.Request.Header.Get("XUserID"), 10, 64)
	user, err := h.repo.GetByID(ctx, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("couldn't find user by mail: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get general account info
// @Tags accounts
// @Accept  json
// @Param Set-Cookie header string true "Authorization token cookie"
// @Produce  json
// @Success 200 {object} dto.UserModelResponseDTO
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /users/:id [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	c.Request.ParseForm()
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("wrong user id: %s", err.Error())})
		return
	}
	user, err := h.repo.GetByID(ctx, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("couldn't find user by mail: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Get user rooms
// @Tags accounts
// @Accept  json
// @Param Set-Cookie header string true "Authorization token cookie"
// @Produce  json
// @Success 200 {object} []models.Room
// @Failure 400 {object} ErrorMessage
// @Failure 401 {object} ErrorMessage
// @Failure 500 {object} ErrorMessage
// @Router /users/:id/rooms [get]
func (h *UserHandler) GetUserRooms(c *gin.Context) {
	ctx := c.Request.Context()
	c.Request.ParseForm()
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("wrong user id: %s", err.Error())})
		return
	}
	rooms, err := h.repo.GetUserRooms(ctx, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("couldn't find user by mail: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, rooms)
}
