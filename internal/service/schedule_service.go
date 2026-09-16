package service

import (
	"errors"
	"time"

	"go-appointment-booking/internal/model"

	"gorm.io/gorm"
)

type ScheduleService struct {
	db *gorm.DB
}

func NewScheduleService(db *gorm.DB) *ScheduleService {
	return &ScheduleService{
		db: db,
	}
}

func (s *ScheduleService) CreateSchedule(schedule *model.DoctorSchedule) error {
	// Check doctor exists
	var doctor model.Doctor

	if err := s.db.First(&doctor, schedule.DoctorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("doctor not found")
		}

		return err
	}

	// Validate time format
	start, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return errors.New("invalid start_time, expected HH:mm")
	}

	end, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return errors.New("invalid end_time, expected HH:mm")
	}

	if !start.Before(end) {
		return errors.New("start_time must be before end_time")
	}

	// Validate break
	if schedule.BreakStartTime != nil || schedule.BreakEndTime != nil {
		if schedule.BreakStartTime == nil || schedule.BreakEndTime == nil {
			return errors.New("break_start_time and break_end_time must be provided together")
		}

		breakStart, err := time.Parse("15:04", *schedule.BreakStartTime)
		if err != nil {
			return errors.New("invalid break_start_time, expected HH:mm")
		}

		breakEnd, err := time.Parse("15:04", *schedule.BreakEndTime)
		if err != nil {
			return errors.New("invalid break_end_time, expected HH:mm")
		}

		if !breakStart.Before(breakEnd) {
			return errors.New("break_start_time must be before break_end_time")
		}

		if breakStart.Before(start) || breakEnd.After(end) {
			return errors.New("break time must be within working hours")
		}
	}

	// Prevent duplicate schedule for same doctor/day
	var existing model.DoctorSchedule

	err = s.db.
		Where("doctor_id = ? AND day_of_week = ?", schedule.DoctorID, schedule.DayOfWeek).
		First(&existing).Error

	if err == nil {
		return errors.New("schedule already exists for this doctor and day")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Create(schedule).Error
}

func (s *ScheduleService) GetSchedules(doctorID uint) ([]model.DoctorSchedule, error) {
	var schedules []model.DoctorSchedule

	if err := s.db.
		Where("doctor_id = ?", doctorID).
		Order("day_of_week ASC").
		Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}
