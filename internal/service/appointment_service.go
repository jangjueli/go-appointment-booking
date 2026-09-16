package service

import (
	"errors"
	"log"
	"time"

	"go-appointment-booking/internal/model"

	"gorm.io/gorm"
)

type AppointmentService struct {
	db *gorm.DB
}

func NewAppointmentService(db *gorm.DB) *AppointmentService {
	return &AppointmentService{
		db: db,
	}
}

type CreateAppointmentInput struct {
	PatientID         uint
	DoctorID          uint
	DepartmentID      uint
	AppointmentTypeID uint
	AppointmentDate   string
	StartTime         string
	Reason            *string
	CreatedBy         string
}

func (s *AppointmentService) CreateAppointment(input CreateAppointmentInput) (*model.Appointment, error) {

	// 1. Parse date
	appointmentDate, err := time.Parse("2006-01-02", input.AppointmentDate)
	if err != nil {
		return nil, errors.New("invalid appointment_date, expected YYYY-MM-DD")
	}

	// 2. Parse start time
	startTime, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return nil, errors.New("invalid start_time, expected HH:mm")
	}

	// 3. Check patient
	var patient model.Patient

	if err := s.db.First(&patient, input.PatientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}

		return nil, err
	}

	// 4. Check doctor
	var doctor model.Doctor

	if err := s.db.First(&doctor, input.DoctorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor not found")
		}

		return nil, err
	}

	if doctor.DepartmentID != input.DepartmentID {
		return nil, errors.New("doctor does not belong to department")
	}

	// 5. Check department
	var department model.Department

	if err := s.db.First(&department, input.DepartmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}

		return nil, err
	}

	// 6. Check appointment type
	var appointmentType model.AppointmentType

	if err := s.db.First(&appointmentType, input.AppointmentTypeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment type not found")
		}

		return nil, err
	}

	// 7. Calculate end time from appointment type duration
	duration := time.Duration(appointmentType.DurationMinutes) * time.Minute
	endTime := startTime.Add(duration)

	// 8. Check doctor schedule
	dayOfWeek := int(appointmentDate.Weekday())

	var schedule model.DoctorSchedule

	err = s.db.
		Where(
			"doctor_id = ? AND day_of_week = ? AND accept_booking = ?",
			input.DoctorID,
			dayOfWeek,
			true,
		).
		First(&schedule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor is not available for booking on this day")
		}

		return nil, err
	}

	scheduleStart, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return nil, errors.New("invalid doctor schedule start time")
	}

	scheduleEnd, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return nil, errors.New("invalid doctor schedule end time")
	}

	// 9. Check working hours
	if startTime.Before(scheduleStart) || endTime.After(scheduleEnd) {
		return nil, errors.New("appointment is outside doctor working hours")
	}

	// 10. Check break time
	if schedule.BreakStartTime != nil && schedule.BreakEndTime != nil {
		breakStart, err := time.Parse("15:04", *schedule.BreakStartTime)
		if err != nil {
			return nil, errors.New("invalid break start time")
		}

		breakEnd, err := time.Parse("15:04", *schedule.BreakEndTime)
		if err != nil {
			return nil, errors.New("invalid break end time")
		}

		if startTime.Before(breakEnd) && endTime.After(breakStart) {
			return nil, errors.New("appointment overlaps doctor break time")
		}
	}

	// 11. Check existing appointments
	var appointments []model.Appointment

	if err := s.db.
		Where(
			"doctor_id = ? AND appointment_date = ? AND status = ?",
			input.DoctorID,
			appointmentDate.Format("2006-01-02"),
			model.AppointmentStatusBooked,
		).
		Find(&appointments).Error; err != nil {
		return nil, err
	}

	for _, appointment := range appointments {
		if IsTimeOverlapping(
			input.StartTime,
			endTime.Format("15:04"),
			appointment.StartTime,
			appointment.EndTime,
		) {
			return nil, errors.New("appointment overlaps existing appointment")
		}
	}

	// 12. Create appointment
	appointment := &model.Appointment{
		PatientID:         input.PatientID,
		DoctorID:          input.DoctorID,
		DepartmentID:      input.DepartmentID,
		AppointmentTypeID: input.AppointmentTypeID,
		AppointmentDate:   appointmentDate,
		StartTime:         input.StartTime,
		EndTime:           endTime.Format("15:04"),
		Status:            model.AppointmentStatusBooked,
		Reason:            input.Reason,
		CreatedBy:         input.CreatedBy,
	}

	if err := s.db.Create(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func (s *AppointmentService) FindAvailableSlots(doctorID uint, date time.Time, appointmentTypeID uint) ([]Slot, error) {

	// 1. Get appointment type
	var appointmentType model.AppointmentType

	if err := s.db.First(&appointmentType, appointmentTypeID).Error; err != nil {
		return nil, err
	}

	// 2. Get doctor's schedule for this day
	dayOfWeek := int(date.Weekday())

	log.Printf("DEBUG date=%v weekday=%d", date, dayOfWeek)

	var schedule model.DoctorSchedule

	err := s.db.
		Where(
			"doctor_id = ? AND day_of_week = ? AND accept_booking = ?",
			doctorID,
			dayOfWeek,
			true,
		).
		First(&schedule).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []Slot{}, nil
		}

		return nil, err
	}

	// DEBUG: ดูค่าที่อ่านจาก database
	log.Printf(
		"DEBUG schedule: start=%v end=%v breakStart=%v breakEnd=%v",
		schedule.StartTime,
		schedule.EndTime,
		schedule.BreakStartTime,
		schedule.BreakEndTime,
	)

	// 3. Generate all slots from working hours
	slots, err := GenerateSlots(
		schedule.StartTime,
		schedule.EndTime,
		appointmentType.DurationMinutes,
	)
	if err != nil {
		return nil, err
	}

	// DEBUG: ดู slots ที่ Generate ออกมา
	log.Printf("DEBUG generated slots: %+v", slots)

	// 4. Remove break-time slots
	slots = FilterBreakSlots(
		slots,
		schedule.BreakStartTime,
		schedule.BreakEndTime,
	)

	// 5. Get booked appointments for this doctor/date
	var appointments []model.Appointment

	if err := s.db.
		Where(
			"doctor_id = ? AND appointment_date = ? AND status = ?",
			doctorID,
			date.Format("2006-01-02"),
			model.AppointmentStatusBooked,
		).
		Find(&appointments).
		Error; err != nil {
		return nil, err
	}

	// 6. Remove booked slots
	slots = FilterBookedSlots(slots, appointments)

	return slots, nil
}
