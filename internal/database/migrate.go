package database

import (
	"go-appointment-booking/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Department{},
		&model.Doctor{},
		&model.Patient{},
		&model.AppointmentType{},
		&model.DoctorSchedule{},
		&model.Appointment{},
	)
}
