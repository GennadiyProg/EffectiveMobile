package request

import (
	"time"

	"github.com/google/uuid"
)

type SubCreateRequest struct {
	ServiceTitle string  `json:"service_title"`
	Price        int     `json:"price"`
	User         string  `json:"user"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date"`
}

type SubUpdateRequest struct {
	ServiceTitle *string `json:"service_title"`
	Price        *int    `json:"price"`
	StartDate    *string `json:"start_date"`
	EndDate      *string `json:"end_date"`
}

type SubUpdate struct {
	ServiceTitle *string
	Price        *int
	StartDate    *time.Time
	EndDate      *time.Time
}

type SubFilterRequest struct {
	StartPeriod  string     `schema:"start_period"`
	EndPeriod    string     `schema:"end_period"`
	UserId       *uuid.UUID `schema:"user_id"`
	ServiceTitle *string    `schema:"service_title"`
}

type SubFilter struct {
	StartPeriod  time.Time
	EndPeriod    time.Time
	UserId       *uuid.UUID
	ServiceTitle *string
}
