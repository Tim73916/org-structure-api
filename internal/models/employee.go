package models

import (
	"errors"
	"strings"
	"time"
)

type Employee struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	DepartmentID uint       `gorm:"not null;index" json:"department_id"`
	FullName     string     `gorm:"not null;size:200" json:"full_name"`
	Position     string     `gorm:"not null;size:200" json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (e *Employee) Validate() error {
	e.FullName = strings.TrimSpace(e.FullName)
	e.Position = strings.TrimSpace(e.Position)

	if e.FullName == "" {
		return errors.New("full_name cannot be empty")
	}

	if len(e.FullName) > 200 {
		return errors.New("full_name must be less than 200 characters")
	}

	if e.Position == "" {
		return errors.New("position cannot be empty")
	}

	if len(e.Position) > 200 {
		return errors.New("position must be less than 200 characters")
	}

	return nil
}
