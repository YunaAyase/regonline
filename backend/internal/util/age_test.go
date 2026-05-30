package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAge(t *testing.T) {
	today := time.Now()

	tests := []struct {
		name      string
		birthDate time.Time
		wantAge   int
	}{
		{
			name:      "exact birthday today",
			birthDate: time.Date(today.Year()-10, today.Month(), today.Day(), 0, 0, 0, 0, time.UTC),
			wantAge:   10,
		},
		{
			name:      "birthday yesterday",
			birthDate: time.Date(today.Year()-10, today.Month(), today.Day()-1, 0, 0, 0, 0, time.UTC),
			wantAge:   10,
		},
		{
			name:      "birthday tomorrow",
			birthDate: time.Date(today.Year()-10, today.Month(), today.Day()+1, 0, 0, 0, 0, time.UTC),
			wantAge:   9,
		},
		{
			name:      "newborn",
			birthDate: today,
			wantAge:   0,
		},
		{
			name:      "exactly 5 years old",
			birthDate: time.Date(today.Year()-5, today.Month(), today.Day(), 0, 0, 0, 0, time.UTC),
			wantAge:   5,
		},
		{
			name:      "exactly 16 years old",
			birthDate: time.Date(today.Year()-16, today.Month(), today.Day(), 0, 0, 0, 0, time.UTC),
			wantAge:   16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.birthDate)
			assert.Equal(t, tt.wantAge, got)
		})
	}
}

func TestCalculateDateRangeByAge(t *testing.T) {
	minAge := 7
	maxAge := 16

	minDate, maxDate := CalculateDateRangeByAge(minAge, maxAge)

	today := time.Now()
	expectedMin := time.Date(today.Year()-maxAge, today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	expectedMax := time.Date(today.Year()-minAge, today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	assert.Equal(t, expectedMin, minDate)
	assert.Equal(t, expectedMax, maxDate)
	assert.True(t, minDate.Before(maxDate))
}
