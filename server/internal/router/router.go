package router

import (
	"github.com/Coderovshik/chat_server/internal/room"
	"github.com/Coderovshik/chat_server/internal/user"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(userHandler *user.Handler, roomHandler *room.Handler) *Router {
	e := gin.Default()

	userGroup := e.Group("/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.POST("/login", userHandler.Login)
		userGroup.DELETE("/logout", userHandler.Logout)
	}

	roomGroup := e.Group("/rooms")
	{
		roomGroup.GET("/", roomHandler.GetRooms)
		roomGroup.POST("/", roomHandler.CreateRoom)
		roomGroup.DELETE("/:id", roomHandler.DeleteRoom)

		roomGroup.GET("/:id", roomHandler.JoinRoom)

		roomGroup.GET("/:id/clients", roomHandler.GetClients)
	}

	return &Router{
		engine: e,
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
