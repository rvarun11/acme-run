package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/services"
	"github.com/google/uuid"

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

func (h *HTTPHandler) ConnectHRM(ctx *gin.Context) {

	type tempDTO struct {
		HRMId   uuid.UUID `json:"hrm_id"`
		Connect bool      `json:"connect"`
	}

	var tempDTOInstance tempDTO
	if err := ctx.ShouldBindJSON(&tempDTOInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	if tempDTOInstance.Connect {
		h.svc.ConnectHRM(tempDTOInstance.HRMId)
	} else {
		h.svc.DisconnectHRM(tempDTOInstance.HRMId)
	}
}
