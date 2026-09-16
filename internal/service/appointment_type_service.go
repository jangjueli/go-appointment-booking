package service

import (
	"go-appointment-booking/internal/model"

	"gorm.io/gorm"
)

type AppointmentTypeService struct {
	db *gorm.DB
}

func NewAppointmentTypeService(db *gorm.DB) *AppointmentTypeService {
	return &AppointmentTypeService{
		db: db,
	}
}

func (s *AppointmentTypeService) GetAppointmentTypes() ([]model.AppointmentType, error) {
	var types []model.AppointmentType

	if err := s.db.
		Order("id ASC").
		Find(&types).Error; err != nil {
		return nil, err
	}

	return types, nil
}
