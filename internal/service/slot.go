package service

import (
	"go-appointment-booking/internal/model"
	"time"
)

type Slot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func GenerateSlots(startTime string, endTime string, durationMinutes int) ([]Slot, error) {
	const timeLayout = "15:04"

	start, err := time.Parse(timeLayout, startTime)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(timeLayout, endTime)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(durationMinutes) * time.Minute

	var slots []Slot

	for current := start; !current.Add(duration).After(end); current = current.Add(duration) {
		slotEnd := current.Add(duration)

		slots = append(slots, Slot{
			StartTime: current.Format(timeLayout),
			EndTime:   slotEnd.Format(timeLayout),
		})
	}

	return slots, nil
}

func IsTimeOverlapping(start string, end string, otherStart string, otherEnd string) bool {
	const timeLayout = "15:04"

	startTime, err := time.Parse(timeLayout, start)
	if err != nil {
		return false
	}

	endTime, err := time.Parse(timeLayout, end)
	if err != nil {
		return false
	}

	otherStartTime, err := time.Parse(timeLayout, otherStart)
	if err != nil {
		return false
	}

	otherEndTime, err := time.Parse(timeLayout, otherEnd)
	if err != nil {
		return false
	}

	return startTime.Before(otherEndTime) &&
		endTime.After(otherStartTime)
}

func FilterBreakSlots(slots []Slot, breakStart *string, breakEnd *string) []Slot {
	if breakStart == nil || breakEnd == nil {
		return slots
	}

	filtered := make([]Slot, 0, len(slots))

	for _, slot := range slots {
		if !IsTimeOverlapping(
			slot.StartTime,
			slot.EndTime,
			*breakStart,
			*breakEnd,
		) {
			filtered = append(filtered, slot)
		}
	}

	return filtered
}

func FilterBookedSlots(slots []Slot, appointments []model.Appointment) []Slot {
	filtered := make([]Slot, 0, len(slots))

	for _, slot := range slots {
		booked := false

		for _, appointment := range appointments {
			if IsTimeOverlapping(
				slot.StartTime,
				slot.EndTime,
				appointment.StartTime,
				appointment.EndTime,
			) {
				booked = true
				break
			}
		}

		if !booked {
			filtered = append(filtered, slot)
		}
	}

	return filtered
}
