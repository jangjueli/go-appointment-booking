package handler

import (
	"net/http"
	"strconv"
	"time"

	"go-appointment-booking/internal/model"
	"go-appointment-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service *service.AppointmentService
}

func NewAppointmentHandler(service *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		service: service,
	}
}

type CreateAppointmentRequest struct {
	PatientID         uint    `json:"patient_id" binding:"required"`
	DoctorID          uint    `json:"doctor_id" binding:"required"`
	DepartmentID      uint    `json:"department_id" binding:"required"`
	AppointmentTypeID uint    `json:"appointment_type_id" binding:"required"`
	AppointmentDate   string  `json:"appointment_date" binding:"required"`
	StartTime         string  `json:"start_time" binding:"required"`
	Reason            *string `json:"reason"`
	CreatedBy         string  `json:"created_by" binding:"required"`
}

type AppointmentResponse struct {
	ID                uint                    `json:"id"`
	PatientID         uint                    `json:"patient_id"`
	DoctorID          uint                    `json:"doctor_id"`
	DepartmentID      uint                    `json:"department_id"`
	AppointmentTypeID uint                    `json:"appointment_type_id"`
	AppointmentDate   string                  `json:"appointment_date"`
	StartTime         string                  `json:"start_time"`
	EndTime           string                  `json:"end_time"`
	Status            model.AppointmentStatus `json:"status"`
	Reason            *string                 `json:"reason"`
	CreatedBy         string                  `json:"created_by"`
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	input := service.CreateAppointmentInput{
		PatientID:         req.PatientID,
		DoctorID:          req.DoctorID,
		DepartmentID:      req.DepartmentID,
		AppointmentTypeID: req.AppointmentTypeID,
		AppointmentDate:   req.AppointmentDate,
		StartTime:         req.StartTime,
		Reason:            req.Reason,
		CreatedBy:         req.CreatedBy,
	}

	appointment, err := h.service.CreateAppointment(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := AppointmentResponse{
		ID:                appointment.ID,
		PatientID:         appointment.PatientID,
		DoctorID:          appointment.DoctorID,
		DepartmentID:      appointment.DepartmentID,
		AppointmentTypeID: appointment.AppointmentTypeID,
		AppointmentDate:   appointment.AppointmentDate.Format("2006-01-02"),
		StartTime:         appointment.StartTime,
		EndTime:           appointment.EndTime,
		Status:            appointment.Status,
		Reason:            appointment.Reason,
		CreatedBy:         appointment.CreatedBy,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AppointmentHandler) FindAvailableSlots(c *gin.Context) {
	doctorID, err := strconv.ParseUint(
		c.Param("doctorID"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid doctor_id",
		})
		return
	}

	dateString := c.Query("date")

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date, expected YYYY-MM-DD",
		})
		return
	}

	appointmentTypeID, err := strconv.ParseUint(
		c.Query("appointment_type_id"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid appointment_type_id",
		})
		return
	}

	slots, err := h.service.FindAvailableSlots(
		uint(doctorID),
		date,
		uint(appointmentTypeID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"doctor_id":           doctorID,
		"date":                dateString,
		"appointment_type_id": appointmentTypeID,
		"slots":               slots,
	})
}
