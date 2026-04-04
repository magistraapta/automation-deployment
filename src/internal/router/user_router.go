package router

import (
	"src/internal/handler"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine, userHandler *handler.UserHandler) {
	users := r.Group("/users")
	{
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
