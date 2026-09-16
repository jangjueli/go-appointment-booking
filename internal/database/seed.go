package database

import (
	"log"
	"time"

	"go-appointment-booking/internal/model"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	// Create department
	department := model.Department{
		Name: "Cardiology",
	}
	log.Println("seeding department")
	if err := db.FirstOrCreate(
		&department,
		model.Department{
			Name: department.Name,
		},
	).Error; err != nil {
		return err
	}

	// Create doctors
	doctors := []model.Doctor{
		{
			Name:         "Dr. Smith",
			DepartmentID: department.ID,
		},
		{
			Name:         "Dr. Johnson",
			DepartmentID: department.ID,
		},
	}
	log.Println("seeding doctors")
	for i := range doctors {
		if err := db.FirstOrCreate(
			&doctors[i],
			model.Doctor{
				Name: doctors[i].Name,
			},
		).Error; err != nil {
			return err
		}
	}

	// Appointment types
	appointmentTypes := []model.AppointmentType{
		{
			Name:            "New patient visit",
			DurationMinutes: 60,
		},
		{
			Name:            "Follow-up visit",
			DurationMinutes: 30,
		},
		{
			Name:            "Consultation",
			DurationMinutes: 30,
		},
		{
			Name:            "Procedure",
			DurationMinutes: 60,
		},
	}
	log.Println("seeding appointment types")
	for i := range appointmentTypes {
		if err := db.FirstOrCreate(
			&appointmentTypes[i],
			model.AppointmentType{
				Name: appointmentTypes[i].Name,
			},
		).Error; err != nil {
			return err
		}
	}

	// Doctor schedules
	schedules := []model.DoctorSchedule{
		{
			DoctorID:       doctors[0].ID,
			DayOfWeek:      1, // Monday
			StartTime:      "09:00",
			EndTime:        "17:00",
			BreakStartTime: stringPtr("12:00"),
			BreakEndTime:   stringPtr("13:00"),
			AcceptBooking:  true,
		},
		{
			DoctorID:       doctors[0].ID,
			DayOfWeek:      2, // Tuesday
			StartTime:      "09:00",
			EndTime:        "17:00",
			BreakStartTime: stringPtr("12:00"),
			BreakEndTime:   stringPtr("13:00"),
			AcceptBooking:  true,
		},
	}
	log.Println("seeding schedules")
	for i := range schedules {
		if err := db.
			Where(
				"doctor_id = ? AND day_of_week = ?",
				schedules[i].DoctorID,
				schedules[i].DayOfWeek,
			).
			FirstOrCreate(&schedules[i]).
			Error; err != nil {
			return err
		}
	}

	// Patient
	patient := model.Patient{
		Name: "John Doe",
	}
	log.Println("seeding patient")
	if err := db.FirstOrCreate(
		&patient,
		model.Patient{
			Name: patient.Name,
		},
	).Error; err != nil {
		return err
	}

	// Sample appointment
	appointment := model.Appointment{
		PatientID:         patient.ID,
		DoctorID:          doctors[0].ID,
		DepartmentID:      department.ID,
		AppointmentTypeID: appointmentTypes[1].ID, // Follow-up visit
		AppointmentDate:   nextMonday(),
		StartTime:         "10:00",
		EndTime:           "10:30",
		Status:            model.AppointmentStatusBooked,
		CreatedBy:         "seed",
	}
	log.Println("seeding appointment")
	if err := db.FirstOrCreate(
		&appointment,
		model.Appointment{
			DoctorID:        appointment.DoctorID,
			AppointmentDate: appointment.AppointmentDate,
			StartTime:       appointment.StartTime,
		},
	).Error; err != nil {
		return err
	}

	log.Println("database seed completed")

	return nil
}

func parseTime(value string) time.Time {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	t, err := time.ParseInLocation("15:04", value, loc)
	if err != nil {
		panic(err)
	}

	return t
}

func timePtr(value time.Time) *time.Time {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func nextMonday() time.Time {
	now := time.Now()
	daysUntilMonday := (8 - int(now.Weekday())) % 7

	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}

	next := now.AddDate(0, 0, daysUntilMonday)

	return time.Date(
		next.Year(),
		next.Month(),
		next.Day(),
		0, 0, 0, 0,
		next.Location(),
	)
}
