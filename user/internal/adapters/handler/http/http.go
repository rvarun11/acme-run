package http

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	gin *gin.Engine
	svc *services.PlayerService
}

func NewPlayerHandler(gin *gin.Engine, playerSvc *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		gin: gin,
		svc: playerSvc,
	}
}

func (h *PlayerHandler) InitRouter() {
	router := h.gin.Group("/api/v1")
	router.GET("/players", h.listPlayers)
	router.POST("/players", h.registerPlayer)
	router.GET("/players/:id", h.getPlayer)
	router.PUT("/players", h.updatePlayer)
	router.DELETE("/players:id", h.deletePlayer)
}

// Players

// @Summary	List Players
// @Tags		players
// @ID			list-players
// @Produce	json
// @Success	200	{array}	http.playerDTO
// @Router		/api/v1/players [get]
func (h *PlayerHandler) listPlayers(ctx *gin.Context) {
	players, err := h.svc.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while fetching players",
		})
		return
	}
	ctx.JSON(http.StatusOK, players)
}

// @Summary	Create a Player
// @Tags		players
// @ID			create-player
// @Produce	json
// @Param		player	body		http.playerDTO	true	"Player data"
// @Success	200		{object}	http.playerDTO
// @Router		/api/v1/players [post]
func (h *PlayerHandler) registerPlayer(ctx *gin.Context) {
	var req playerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	player, err := h.svc.Register(req.toAggregate())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to register, something occured",
		})
		return
	}

	res := fromAggregate(player)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "New player created successfully",
		"player":  res,
	})
}

// @Summary	Get Player by ID
// @Tags		players
// @ID			get-player
// @Produce	json
// @Param		id	path		string	true	"Player ID (UUID)"
// @Success	200	{object}	http.playerDTO
// @Router		/api/v1/player/{id} [get]
func (h *PlayerHandler) getPlayer(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	p, err := h.svc.Get(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while fetching player",
		})
		return
	}
	player := fromAggregate(p)
	ctx.JSON(http.StatusOK, player)
}

// @Summary	Update Player
// @Tags		players
// @ID			update-challenge
// @Produce	json
// @Param		player	body	http.playerDTO	true	"Player data"
// @Success	204
// @Router		/api/v1/player [put]
func (h *PlayerHandler) updatePlayer(ctx *gin.Context) {
	var req playerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	p, err := h.svc.Update(req.toAggregate())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while updating player",
		})
		return
	}
	res := fromAggregate(p)
	ctx.JSON(http.StatusNoContent, gin.H{
		// 204 returns no body so this can be changed to 200 if body is needed
		"player":  res,
		"message": "Player updated successfully",
	})
}

// @Summary	Delete a Player by ID
// @Tags		players
// @ID			delete-player
// @Produce	json
// @Param		id	path	string	true	"Player ID (UUID) to be deleted"
// @Success	204	"No Content"
// @Router		/api/v1/players/{id} [delete]
func (h *PlayerHandler) deletePlayer(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	// Delete the Challenge
	err = h.svc.DeletePlayerByID(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while deleting the player",
		})
		return
	}
	ctx.JSON(http.StatusNoContent, err)
}
