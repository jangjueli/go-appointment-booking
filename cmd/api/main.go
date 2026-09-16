package main

import (
	"log"

	"go-appointment-booking/internal/database"
	"go-appointment-booking/internal/handler"
	"go-appointment-booking/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// db
	db, err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	log.Println("database connected:", db != nil)
	// migrate
	if err := database.Migrate(db); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("database migrated successfully")
	// seed data
	if err := database.Seed(db); err != nil {
		log.Fatal("failed to seed database:", err)
	}
	log.Println("database seed completed successfully")

	scheduleService := service.NewScheduleService(db)
	appointmentService := service.NewAppointmentService(db)
	appointmentTypeService := service.NewAppointmentTypeService(db)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	appointmentTypeHandler := handler.NewAppointmentTypeHandler(appointmentTypeService)

	router := gin.Default()
	router.POST("/doctors/:doctorID/schedules", scheduleHandler.CreateSchedule)
	router.GET("/doctors/:doctorID/schedules", scheduleHandler.GetSchedules)
	router.GET("/doctors/:doctorID/available-slots", appointmentHandler.FindAvailableSlots)
	router.POST("/appointments", appointmentHandler.CreateAppointment)
	router.GET("/appointment-types", appointmentTypeHandler.GetAppointmentTypes)

	log.Println("server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
