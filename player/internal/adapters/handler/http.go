package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvs_/player/internal/core/services"

	"github.com/CAS735-F23/macrun-teamvs_/player/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	svc services.PlayerService
}

func NewHTTPHandler(PlayerService services.PlayerService) *HTTPHandler {
	return &HTTPHandler{
		svc: PlayerService,
	}
}

func (h *HTTPHandler) ListPlayers(ctx *gin.Context) {

	players, err := h.svc.List()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, players)

}

func (h *HTTPHandler) CreatePlayer(ctx *gin.Context) {
	var p domain.Player
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	// TODO: Should this be here or in the service?
	// Keeping it here for now so that uuid can be sent back
	player, err := domain.NewPlayer(p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	// TODO: The two error handling are the same, it can be refactored
	err = h.svc.Create(player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":      player.User.ID,
		"message": "New player created successfully",
	})
}

func (h *HTTPHandler) GetPlayer(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	player, err := h.svc.Get(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, player)
}

// TODO: To be implemented
func (h *HTTPHandler) UpdatePlayer(ctx *gin.Context) {

}
