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
	router.DELETE("/challenges/:id", h.deleteChallengeById)

	// Badges
	router.GET("/badges", h.listBadgesByPlayerID)
}

// Challenges

// @Summary	Create a Challenge
// @Tags		challenges
// @ID			create-challenge
// @Produce	json
// @Param		challenge	body		http.challengeDTO	true	"Challenge data"
// @Success	200			{object}	http.challengeDTO
// @Router		/api/v1/challenges [post]
func (h *ChallengeHandler) createChallenge(ctx *gin.Context) {
	var req challengeDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	ch, err := h.svc.CreateChallenge(toAggregate(&req))
	res := fromAggregate(ch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to create challenge",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"challenge": res,
		"message":   "New challenge created successfully",
	})
}

// @Summary	Get Challenge by ID
// @Tags		challenges
// @ID			get-challenge
// @Produce	json
// @Param		id	path		string	true	"Challenge ID (UUID)"
// @Success	200	{object}	http.challengeDTO
// @Router		/api/v1/challenges/{id} [get]
func (h *ChallengeHandler) getChallengeByID(ctx *gin.Context) {
	cid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	ch, err := h.svc.GetChallengeByID(cid)
	res := fromAggregate(ch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while fetching the challenge",
		})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary	Update Challenge
// @Tags		challenges
// @ID			update-challenge
// @Produce	json
// @Param		challenge	body	http.challengeDTO	true	"Challenge data"
// @Success	204
// @Router		/api/v1/challenges [put]
func (h *ChallengeHandler) updateChallenge(ctx *gin.Context) {
	var req *challengeDTO
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	ch, err := h.svc.UpdateChallenge(toAggregate(req))
	res := fromAggregate(ch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while updating challenge",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		// 204 returns no body so this can be changed to 200 if body is needed/
		"player":  res,
		"message": "Player updated successfully",
	})
}

// @Summary	List Challenges
// @Tags		challenges
// @ID			list-challenges
// @Produce	json
// @Success	200	{array}	http.challengeDTO
// @Router		/api/v1/challenges [get]
func (h *ChallengeHandler) listChallenges(ctx *gin.Context) {
	status := ctx.Query("status")
	chs, err := h.svc.ListChallenges(status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to fetch list of challenges",
		})
		return
	}
	ctx.JSON(http.StatusOK, chs)
}

// @Summary	Delete a Challenge by ID
// @Tags		challenges
// @ID			delete-challenge
// @Produce	json
// @Param		id	path	string	true	"Challenge ID (UUID) to be deleted"
// @Success	204	"No Content"
// @Router		/api/v1/challenges/{id} [delete]
func (h *ChallengeHandler) deleteChallengeById(ctx *gin.Context) {
	cid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	// Delete the Challenge
	err = h.svc.DeleteChallengeByID(cid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while deleting the challenge",
		})
		return
	}
	ctx.JSON(http.StatusNoContent, err)
}

// Badges

// @Summary	List Badges by Player ID
// @ID			list-badges
// @Tags		badges
// @Produce	json
// @Param		player_id	query	string	true	"Player ID (UUID)"
// @Success	200			{array}	domain.Badge
// @Router		/api/v1/badges [get]
func (h *ChallengeHandler) listBadgesByPlayerID(ctx *gin.Context) {
	pid, err := uuid.Parse(ctx.Query("player_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	badges, err := h.svc.ListBadgesByPlayerID(pid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error occured while fetching badges for player",
		})
		return
	}
	ctx.JSON(http.StatusOK, badges)
}
