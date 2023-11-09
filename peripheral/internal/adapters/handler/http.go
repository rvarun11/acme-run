package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc services.PeripheralService
}

func NewHTTPHandler(PeripheralService services.PeripheralService) *HTTPHandler {
	return &HTTPHandler{
		svc: PeripheralService,
	}
}

func (h *HTTPHandler) ConnectPeripheral(ctx *gin.Context) {

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
		h.svc.ConnectPeripheral(tempDTOInstance.HRMId)
	} else {
		h.svc.DisconnectPeripheral(tempDTOInstance.HRMId)
	}
}
