package handler

import (
	"net/http"

	"go-appointment-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type AppointmentTypeHandler struct {
	service *service.AppointmentTypeService
}

func NewAppointmentTypeHandler(service *service.AppointmentTypeService) *AppointmentTypeHandler {
	return &AppointmentTypeHandler{
		service: service,
	}
}

func (h *AppointmentTypeHandler) GetAppointmentTypes(c *gin.Context) {
	types, err := h.service.GetAppointmentTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types)
}
