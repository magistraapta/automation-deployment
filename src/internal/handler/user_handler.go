package handler

import (
	"src/internal/model"

	"github.com/gin-gonic/gin"
)

var users = []model.User{
	model.User{
		ID:       "1",
		Username: "John Doe",
		Email:    "john.doe@example.com",
	},
	model.User{
		ID:       "2",
		Username: "Jane Doe",
		Email:    "jane.doe@example.com",
	},
}

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	c.JSON(200, users)
}
