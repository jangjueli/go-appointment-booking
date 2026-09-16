package model

import "time"

type AppointmentStatus string

const (
	AppointmentStatusBooked    AppointmentStatus = "booked"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	ID                uint              `gorm:"primaryKey"`
	PatientID         uint              `gorm:"not null"`
	DoctorID          uint              `gorm:"not null"`
	DepartmentID      uint              `gorm:"not null"`
	AppointmentTypeID uint              `gorm:"not null"`
	AppointmentDate   time.Time         `gorm:"type:date;not null"`
	StartTime         string            `gorm:"type:varchar(5);not null"`
	EndTime           string            `gorm:"type:varchar(5);not null"`
	Status            AppointmentStatus `gorm:"type:varchar(20);not null"`
	Reason            *string
	CreatedBy         string `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time

	Patient         Patient
	Doctor          Doctor
	Department      Department
	AppointmentType AppointmentType
}
