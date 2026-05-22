package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepartmentValidate(t *testing.T) {
	tests := []struct {
		name       string
		department Department
		wantError  bool
	}{
		{
			name:       "valid department",
			department: Department{Name: "IT"},
			wantError:  false,
		},
		{
			name:       "empty name",
			department: Department{Name: ""},
			wantError:  true,
		},
		{
			name:       "name too long",
			department: Department{Name: string(make([]byte, 250))},
			wantError:  true,
		},
		{
			name:       "self parent",
			department: Department{ID: 1, Name: "Test", ParentID: func() *uint { id := uint(1); return &id }()},
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.department.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
