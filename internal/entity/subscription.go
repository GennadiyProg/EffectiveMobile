package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID           uuid.UUID  `json:"id"`
	ServiceTitle string     `json:"service_title"`
	Price        int        `json:"price"`
	User         string     `json:"user"`
	StartDate    time.Time  `json:"start_date" gorm:"type:date"`
	EndDate      *time.Time `json:"end_date" gorm:"type:date"`
}
