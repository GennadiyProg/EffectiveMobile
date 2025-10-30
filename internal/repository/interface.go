package repository

import (
	"effective_mobile/internal/dto/request"
	"effective_mobile/internal/entity"

	"github.com/google/uuid"
)

type SubRepository interface {
	Create(sub *entity.Subscription) error
	Get(subId uuid.UUID) (*entity.Subscription, error)
	Update(subId uuid.UUID, subUpdate *request.SubUpdate) (*entity.Subscription, error)
	Delete(subId uuid.UUID) error
	GetPage(pageNumber int, count int) ([]*entity.Subscription, error)
	GetPriceSumByFilter(filter *request.SubFilter) (int, error)
}
