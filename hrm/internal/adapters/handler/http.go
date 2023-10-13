package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/services"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/domain"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc services.HRMService
}

func NewHTTPHandler(HRMService services.HRMService) *HTTPHandler {
	return &HTTPHandler{
		svc: HRMService,
	}
}

func (h *HTTPHandler) ListHRM(ctx *gin.Context) {

	hrms, err := h.svc.List()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, hrms)

}

func (h *HTTPHandler) CreateHRM(ctx *gin.Context) {
	var hrms domain.HRM
	if err := ctx.ShouldBindJSON(&hrms); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	// TODO: Should this be here or in the service?
	// Keeping it here for now so that uuid can be sent back
	hrms, err := domain.NewHRM(hrms)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	// TODO: The two error handling are the same, it can be refactored
	err = h.svc.Create(hrms)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":      hrms.HRMId,
		"message": "New HRM created successfully",
	})
}

func (h *HTTPHandler) GetHRM(ctx *gin.Context) {
	var hid string
	hid = ctx.Param("id")

	hrms, err := h.svc.Get(hid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, hrms)
}

// TODO: To be implemented
func (h *HTTPHandler) UpdateHRM(ctx *gin.Context) {

}
