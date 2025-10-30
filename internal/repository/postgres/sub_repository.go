package postgres

import (
	"effective_mobile/internal/dto/request"
	"effective_mobile/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r SubscriptionRepository) Create(sub *entity.Subscription) error {
	return r.db.Create(&sub).Error
}

func (r SubscriptionRepository) Get(subId uuid.UUID) (*entity.Subscription, error) {
	var sub *entity.Subscription
	result := r.db.Model(&entity.Subscription{}).Where("id = ?", subId).Scan(&sub)
	return sub, result.Error
}

func (r SubscriptionRepository) Update(subId uuid.UUID, subUpdate *request.SubUpdate) (*entity.Subscription, error) {
	fieldUpdate := map[string]interface{}{}
	if subUpdate.Price != nil && *subUpdate.Price > 0 {
		fieldUpdate["price"] = subUpdate.Price
	}
	if subUpdate.ServiceTitle != nil && *subUpdate.ServiceTitle != "" {
		fieldUpdate["service_title"] = subUpdate.ServiceTitle
	}
	if subUpdate.StartDate != nil {
		fieldUpdate["start_date"] = subUpdate.StartDate
	}
	if subUpdate.EndDate != nil {
		fieldUpdate["end_date"] = subUpdate.EndDate
	}
	var updatedSub entity.Subscription
	result := r.db.Model(&entity.Subscription{}).Where("id = ?", subId).Updates(fieldUpdate).Scan(&updatedSub)
	return &updatedSub, result.Error
}

func (r SubscriptionRepository) Delete(subId uuid.UUID) error {
	return r.db.Delete(&entity.Subscription{}, "id = ?", subId).Error
}

func (r SubscriptionRepository) GetPage(pageNumber int, count int) ([]*entity.Subscription, error) {
	var subs []*entity.Subscription
	result := r.db.Model(&entity.Subscription{}).Order("start_date desc").Offset(pageNumber).Limit(count).Find(&subs)
	return subs, result.Error
}

func (r SubscriptionRepository) GetPriceSumByFilter(filter *request.SubFilter) (int, error) {
	query := r.db.Model(&entity.Subscription{}).Select("SUM(price) AS price").
		Where("start_date >= ? and start_date < ?", filter.StartPeriod, filter.EndPeriod)
	if filter.ServiceTitle != nil && *filter.ServiceTitle != "" {
		query = query.Where("service_title = ?", filter.ServiceTitle)
	}
	if filter.UserId != nil && *filter.UserId != uuid.Nil {
		query = query.Where("\"user\" = ?", filter.UserId)
	}
	var count *int
	result := query.Scan(&count)
	if count == nil {
		return 0, nil
	}
	return *count, result.Error
}
