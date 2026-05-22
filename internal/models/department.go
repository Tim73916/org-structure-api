package models

import (
	"errors"
	"strings"
	"time"
)

type Department struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"not null;size:200" json:"name"`
	ParentID  *uint        `gorm:"index" json:"parent_id,omitempty"`
	Parent    *Department  `gorm:"foreignKey:ParentID" json:"-"`
	Children  []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Employees []Employee   `json:"employees,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

func (d *Department) Validate() error {
	d.Name = strings.TrimSpace(d.Name)

	if d.Name == "" {
		return errors.New("name cannot be empty")
	}

	if len(d.Name) > 200 {
		return errors.New("name must be less than 200 characters")
	}

	if d.ParentID != nil && d.ID == *d.ParentID {
		return errors.New("department cannot be parent of itself")
	}

	return nil
}
