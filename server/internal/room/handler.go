package room

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Coderovshik/chat_server/internal/config"
	"github.com/Coderovshik/chat_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub        RoomService
	upgrader   *websocket.Upgrader
	signingKey string
}

func NewHandler(s RoomService, u *websocket.Upgrader, cfg *config.Config) *Handler {
	return &Handler{
		hub:        s,
		upgrader:   u,
		signingKey: cfg.SigningKey,
	}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.hub.CreateRoom(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) JoinRoom(c *gin.Context) {
	roomId := c.Param("id")
	if !h.hub.IsRoomExist(roomId) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("room %s does not exist", roomId),
		})
		return
	}

	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	userClaims, err := util.ParseUserClaims(cookie, h.signingKey)
	if err != nil {
		if errors.Is(err, util.ErrUnknownClaimsType) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to join the room",
		})
		return
	}

	client := NewClient(userClaims.ID, userClaims.Username, roomId, conn)
	h.hub.JoinRoom(roomId, client)
}

func (h *Handler) DeleteRoom(c *gin.Context) {
	roomId := c.Param("id")
	err := h.hub.DeleteRoom(roomId)
	if err != nil {
		if errors.Is(err, ErrRoomNotExist) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("room %s does not exist", roomId),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := &DeleteRoomResponse{
		ID: roomId,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetRooms(c *gin.Context) {
	res, err := h.hub.GetRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetClients(c *gin.Context) {
	roomId := c.Param("id")
	res, err := h.hub.GetClients(roomId)
	if err != nil {
		if errors.Is(err, ErrRoomNotExist) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("room %s does not exist", roomId),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
