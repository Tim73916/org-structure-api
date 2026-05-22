package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployeeValidate(t *testing.T) {
	tests := []struct {
		name      string
		employee  Employee
		wantError bool
	}{
		{
			name:      "valid employee",
			employee:  Employee{FullName: "Ivan Petrov", Position: "Developer"},
			wantError: false,
		},
		{
			name:      "empty full name",
			employee:  Employee{FullName: "", Position: "Developer"},
			wantError: true,
		},
		{
			name:      "empty position",
			employee:  Employee{FullName: "Ivan Petrov", Position: ""},
			wantError: true,
		},
		{
			name:      "full name too long",
			employee:  Employee{FullName: string(make([]byte, 250)), Position: "Developer"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.employee.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
