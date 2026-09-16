package model

import "time"

type AppointmentType struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"not null;unique"`
	DurationMinutes int    `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Appointments []Appointment
}
