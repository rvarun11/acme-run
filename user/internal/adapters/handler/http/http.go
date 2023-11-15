package http

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/services"
	"github.com/CAS735-F23/macrun-teamvsl/user/logger"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
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

	router.POST("/players", h.registerPlayer)
	router.GET("/players/:id", h.getPlayer)
	router.PUT("/players", h.updatePlayer)
	router.GET("/players", h.listPlayers)
}

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
	logger.Info("player registered")
	res := fromAggregate(player)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "New player created successfully",
		"player":  res,
	})
}

func (h *PlayerHandler) getPlayer(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	p, err := h.svc.Get(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	player := fromAggregate(p)
	ctx.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) updatePlayer(ctx *gin.Context) {
	var req playerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	p, err := h.svc.Update(req.toAggregate())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
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

func (h *PlayerHandler) listPlayers(ctx *gin.Context) {
	players, err := h.svc.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, players)
}
