package response

import (
	"effective_mobile/internal/entity"
	"time"

	"github.com/google/uuid"
)

type SubResponse struct {
	ID           *uuid.UUID `json:"id"`
	ServiceTitle *string    `json:"service_title"`
	Price        *int       `json:"price"`
	User         *string    `json:"user"`
	StartDate    *string    `json:"start_date"`
	EndDate      *string    `json:"end_date"`
}

func NewSubResponse(sub *entity.Subscription) *SubResponse {
	if sub == nil {
		return nil
	}
	return &SubResponse{
		ID:           &sub.ID,
		ServiceTitle: &sub.ServiceTitle,
		Price:        &sub.Price,
		User:         &sub.User,
		StartDate:    formDate(&sub.StartDate),
		EndDate:      formDate(sub.EndDate),
	}
}

func formDate(date *time.Time) *string {
	if date == nil {
		return nil
	}
	dateStr := date.Format("2006-01")
	return &dateStr
}

type SubSumResponse struct {
	Result int `json:"result"`
}
