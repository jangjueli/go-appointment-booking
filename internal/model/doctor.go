package model

import "time"

type Doctor struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	DepartmentID uint   `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Department   Department
	Schedules    []DoctorSchedule
	Appointments []Appointment
}
