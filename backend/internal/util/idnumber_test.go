package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeIDNumber(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"411728201503150012", "411728201503150012"},
		{"41172820150315001x", "41172820150315001X"},
		{"  41172820150315001X  ", "41172820150315001X"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormalizeIDNumber(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateIDNumberLength(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"411728201503150012", true},
		{"41172820150315", false},
		{"4117282015031500123", false},
		{"", false},
		{"41172820150315001", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ValidateIDNumberLength(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateBirthDateMatchesID(t *testing.T) {
	tests := []struct {
		name      string
		birthDate string
		idNumber  string
		want      bool
	}{
		{
			name:      "matching date",
			birthDate: "20150315",
			idNumber:  "411728201503150012",
			want:      true,
		},
		{
			name:      "mismatched date",
			birthDate: "20150315",
			idNumber:  "411728202001010014",
			want:      false,
		},
		{
			name:      "short id number",
			birthDate: "20150315",
			idNumber:  "41172820150315",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateBirthDateMatchesID(tt.birthDate, tt.idNumber)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateGenderMatchesID(t *testing.T) {
	tests := []struct {
		name   string
		gender string
		idNum  string
		want   bool
	}{
		{"male match", "男", "411728201503150012", true},
		{"male mismatch", "女", "411728201503150012", false},
		{"female match", "女", "411728201503150024", true},
		{"female mismatch", "男", "411728201503150024", false},
		{"short id", "男", "4117282015031500", false},
		{"invalid gender", "其他", "411728201503150012", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateGenderMatchesID(tt.gender, tt.idNum)
			assert.Equal(t, tt.want, got)
		})
	}
}
