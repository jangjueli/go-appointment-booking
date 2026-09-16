package service

import (
	"testing"
)

func TestGenerateSlots(t *testing.T) {
	slots, err := GenerateSlots("09:00", "12:00", 30)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Slot{
		{StartTime: "09:00", EndTime: "09:30"},
		{StartTime: "09:30", EndTime: "10:00"},
		{StartTime: "10:00", EndTime: "10:30"},
		{StartTime: "10:30", EndTime: "11:00"},
		{StartTime: "11:00", EndTime: "11:30"},
		{StartTime: "11:30", EndTime: "12:00"},
	}

	if len(slots) != len(expected) {
		t.Fatalf("expected %d slots, got %d", len(expected), len(slots))
	}

	for i := range expected {
		if slots[i] != expected[i] {
			t.Errorf(
				"slot %d: expected %+v, got %+v",
				i,
				expected[i],
				slots[i],
			)
		}
	}
}

func TestFilterBreakSlots(t *testing.T) {
	slots := []Slot{
		{StartTime: "09:00", EndTime: "09:30"},
		{StartTime: "09:30", EndTime: "10:00"},
		{StartTime: "10:00", EndTime: "10:30"},
		{StartTime: "10:30", EndTime: "11:00"},
		{StartTime: "11:00", EndTime: "11:30"},
		{StartTime: "11:30", EndTime: "12:00"},
		{StartTime: "12:00", EndTime: "12:30"},
		{StartTime: "12:30", EndTime: "13:00"},
		{StartTime: "13:00", EndTime: "13:30"},
	}

	breakStart := "12:00"
	breakEnd := "13:00"

	filtered := FilterBreakSlots(
		slots,
		&breakStart,
		&breakEnd,
	)

	if len(filtered) != 7 {
		t.Fatalf("expected 7 slots, got %d", len(filtered))
	}

	for _, slot := range filtered {
		if slot.StartTime == "12:00" ||
			slot.StartTime == "12:30" {
			t.Errorf("break slot should not be available: %+v", slot)
		}
	}
}

func TestIsTimeOverlapping(t *testing.T) {
	tests := []struct {
		name       string
		start      string
		end        string
		otherStart string
		otherEnd   string
		expected   bool
	}{
		{
			name:       "overlapping",
			start:      "10:15",
			end:        "10:45",
			otherStart: "10:00",
			otherEnd:   "10:30",
			expected:   true,
		},
		{
			name:       "not overlapping",
			start:      "10:30",
			end:        "11:00",
			otherStart: "10:00",
			otherEnd:   "10:30",
			expected:   false,
		},
		{
			name:       "completely inside",
			start:      "10:10",
			end:        "10:20",
			otherStart: "10:00",
			otherEnd:   "10:30",
			expected:   true,
		},
		{
			name:       "completely outside",
			start:      "11:00",
			end:        "11:30",
			otherStart: "10:00",
			otherEnd:   "10:30",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTimeOverlapping(
				tt.start,
				tt.end,
				tt.otherStart,
				tt.otherEnd,
			)

			if result != tt.expected {
				t.Errorf(
					"expected %v, got %v",
					tt.expected,
					result,
				)
			}
		})
	}
}
