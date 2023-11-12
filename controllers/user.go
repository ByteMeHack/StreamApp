package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/folklinoff/hack_and_change/models"
	auth "github.com/folklinoff/hack_and_change/pkg/auth"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db UserRepository
}

type UserRepository interface {
	Save(ctx context.Context, user models.User) (int64, error)
	GetByMail(ctx context.Context, mail string) (int64, error)
}

type SignUpResponseDTO struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

func (h *UserHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: fmt.Sprintf("failed to bind user model: %s", err.Error())})
		return
	}
	_, err := h.db.GetByMail(ctx, user.Email)
	switch {
	case err == nil:
		c.JSON(http.StatusBadRequest, ErrorMessage{Message: "user already exists"})
		return
	case err != sql.ErrNoRows:
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't get data about user: %s", err.Error())})
		return
	}
	id, err := h.db.Save(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't create user: %s", err.Error())})
		return
	}
	token, err := auth.CreateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage{Message: fmt.Sprintf("couldn't create token: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, SignUpResponseDTO{
		ID:    id,
		Token: token,
	})
}
