package handler

import (
	"net/http"
	"strconv"

	"go-appointment-booking/internal/model"
	"go-appointment-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

type CreateScheduleRequest struct {
	DayOfWeek      int     `json:"day_of_week" binding:"required,min=0,max=6"`
	StartTime      string  `json:"start_time" binding:"required"`
	EndTime        string  `json:"end_time" binding:"required"`
	BreakStartTime *string `json:"break_start_time"`
	BreakEndTime   *string `json:"break_end_time"`
	AcceptBooking  bool    `json:"accept_booking"`
}

type ScheduleResponse struct {
	ID             uint    `json:"id"`
	DoctorID       uint    `json:"doctor_id"`
	DayOfWeek      int     `json:"day_of_week"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	BreakStartTime *string `json:"break_start_time"`
	BreakEndTime   *string `json:"break_end_time"`
	AcceptBooking  bool    `json:"accept_booking"`
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctorID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid doctor_id",
		})
		return
	}

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	schedule := model.DoctorSchedule{
		DoctorID:       uint(doctorID),
		DayOfWeek:      req.DayOfWeek,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		BreakStartTime: req.BreakStartTime,
		BreakEndTime:   req.BreakEndTime,
		AcceptBooking:  req.AcceptBooking,
	}

	if err := h.service.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, ScheduleResponse{
		ID:             schedule.ID,
		DoctorID:       schedule.DoctorID,
		DayOfWeek:      schedule.DayOfWeek,
		StartTime:      schedule.StartTime,
		EndTime:        schedule.EndTime,
		BreakStartTime: schedule.BreakStartTime,
		BreakEndTime:   schedule.BreakEndTime,
		AcceptBooking:  schedule.AcceptBooking,
	})
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctorID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid doctor_id",
		})
		return
	}

	schedules, err := h.service.GetSchedules(uint(doctorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	responses := make([]ScheduleResponse, 0, len(schedules))

	for _, schedule := range schedules {
		responses = append(responses, ScheduleResponse{
			ID:             schedule.ID,
			DoctorID:       schedule.DoctorID,
			DayOfWeek:      schedule.DayOfWeek,
			StartTime:      schedule.StartTime,
			EndTime:        schedule.EndTime,
			BreakStartTime: schedule.BreakStartTime,
			BreakEndTime:   schedule.BreakEndTime,
			AcceptBooking:  schedule.AcceptBooking,
		})
	}

	c.JSON(http.StatusOK, responses)
}
