package model

import "time"

type DoctorSchedule struct {
	ID             uint    `gorm:"primaryKey"`
	DoctorID       uint    `gorm:"not null"`
	DayOfWeek      int     `gorm:"not null"`
	StartTime      string  `gorm:"type:varchar(5);not null"`
	EndTime        string  `gorm:"type:varchar(5);not null"`
	BreakStartTime *string `gorm:"type:varchar(5)"`
	BreakEndTime   *string `gorm:"type:varchar(5)"`
	AcceptBooking  bool    `gorm:"not null;default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Doctor Doctor
}
