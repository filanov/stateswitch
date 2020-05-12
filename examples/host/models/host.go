package models

import (
	"github.com/google/uuid"
)

type Host struct {
	ID     uuid.UUID `json:"id" gorm:"primary_key"`
	Status string    `json:"status"`
	HwInfo *bool     `json:"hw_info"` // probably should be string or json but to make it simple we will work like this
}
