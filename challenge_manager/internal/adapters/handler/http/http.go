package http

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChallengeHandler struct {
	gin *gin.Engine
	svc services.ChallengeService
}

func NewChallengeHandler(gin *gin.Engine, challengeSvc services.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{
		gin: gin,
		svc: challengeSvc,
	}
}

func (h *ChallengeHandler) InitRouter() {
	router := h.gin.Group("/api/v1")
	router.POST("/challenges", h.createChallenge)
	router.GET("/challenges/:id", h.getChallengeByID)
	router.PUT("/challenges", h.updateChallenge)
	router.GET("/challenges", h.listChallenges)
}

func (h *ChallengeHandler) createChallenge(ctx *gin.Context) {
	var req challengeDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ch, err := h.svc.Create(toAggregate(&req))
	res := fromAggregate(ch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"challenge": res,
		"message":   "New challenge created successfully",
	})
}

func (h *ChallengeHandler) getChallengeByID(ctx *gin.Context) {
	cid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ch, err := h.svc.GetByID(cid)
	res := fromAggregate(ch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *ChallengeHandler) updateChallenge(ctx *gin.Context) {
	var req *challengeDTO
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	ch, err := h.svc.Update(toAggregate(req))
	res := fromAggregate(ch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		// 204 returns no body so this can be changed to 200 if body is needed
		"player":  res,
		"message": "Player updated successfully",
	})
}

func (h *ChallengeHandler) listChallenges(ctx *gin.Context) {
	status := ctx.Query("status")
	chs, err := h.svc.List(status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, chs)
}
