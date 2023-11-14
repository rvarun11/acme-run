package http

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChallengeHandler struct {
	gin *gin.Engine
	svc *services.ChallengeService
}

func NewChallengeHandler(gin *gin.Engine, challengeSvc *services.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{
		gin: gin,
		svc: challengeSvc,
	}
}

func (h *ChallengeHandler) InitRouter() {
	router := h.gin.Group("/api/v1")

	// Challenges
	router.POST("/challenges", h.createChallenge)
	router.GET("/challenges/:id", h.getChallengeByID)
	router.PUT("/challenges", h.updateChallenge)
	router.GET("/challenges", h.listChallenges)

	// Badges: This may end up in a separate handler
	router.GET("/badges/:player_id", h.listBadgesByPlayerID)
	// ChallengeStats: This may end up in a separate handler
	router.GET("/stats/:player_id", h.listChallengeStatsByPlayerID)
}

// Challenges

func (h *ChallengeHandler) createChallenge(ctx *gin.Context) {
	var req challengeDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ch, err := h.svc.CreateChallenge(toAggregate(&req))
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
	ch, err := h.svc.GetChallengeByID(cid)
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
			"error": err,
		})
		return
	}

	ch, err := h.svc.UpdateChallenge(toAggregate(req))
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
	chs, err := h.svc.ListChallenges(status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, chs)
}

// Badges
func (h *ChallengeHandler) listBadgesByPlayerID(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("player_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	chs, err := h.svc.ListBadgesByPlayerID(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, chs)
}

// ChallengeStats
func (h *ChallengeHandler) listChallengeStatsByPlayerID(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Param("player_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	csArr, err := h.svc.ListChallengeStatsByPlayerID(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, csArr)
}
