package service

import (
	"effective_mobile/internal/dto/request"
	"effective_mobile/internal/dto/response"
	"effective_mobile/internal/entity"
	"effective_mobile/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	subRepository repository.SubRepository
}

func NewSubscriptionService(subRepository repository.SubRepository) *SubscriptionService {
	return &SubscriptionService{
		subRepository: subRepository,
	}
}

func (service SubscriptionService) CreateSub(sub *request.SubCreateRequest) (*response.SubResponse, error) {
	if sub.Price < 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}
	startDate, err := parseDate(&sub.StartDate)
	if err != nil {
		return nil, fmt.Errorf("%w (start_date)", err)
	} else if startDate == nil {
		return nil, fmt.Errorf("start date must be entered")
	}
	endDate, err := parseDate(sub.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%w (end_date)", err)
	}

	createdSub := entity.Subscription{
		ID:           uuid.New(),
		ServiceTitle: sub.ServiceTitle,
		Price:        sub.Price,
		User:         sub.User,
		StartDate:    *startDate,
		EndDate:      endDate,
	}

	if err := service.subRepository.Create(&createdSub); err != nil {
		return nil, err
	} else {
		return response.NewSubResponse(&createdSub), nil
	}
}

func parseDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}
	if date, err := time.Parse("2006-01", *dateStr); err != nil {
		return nil, fmt.Errorf("could not parse date, must been format 'yyyy-mm'")
	} else {
		return &date, nil
	}
}

func (service SubscriptionService) GetSub(subId uuid.UUID) (*response.SubResponse, error) {
	if sub, err := service.subRepository.Get(subId); err != nil {
		return nil, err
	} else {
		return response.NewSubResponse(sub), nil
	}
}

func (service SubscriptionService) UpdateSub(subId uuid.UUID, subUpdate *request.SubUpdateRequest) (*response.SubResponse, error) {
	if subUpdate.Price != nil && *subUpdate.Price < 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}
	startDate, err := parseDate(subUpdate.StartDate)
	if err != nil {
		return nil, fmt.Errorf("%w (start_date)", err)
	}
	endDate, err := parseDate(subUpdate.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%w (end_date)", err)
	}

	updatedSub := request.SubUpdate{
		ServiceTitle: subUpdate.ServiceTitle,
		Price:        subUpdate.Price,
		StartDate:    startDate,
		EndDate:      endDate,
	}

	if sub, err := service.GetSub(subId); err != nil || sub == nil {
		return nil, err
	}

	if updatedSub, err := service.subRepository.Update(subId, &updatedSub); err != nil {
		return nil, err
	} else {
		return response.NewSubResponse(updatedSub), nil
	}
}

func (service SubscriptionService) DeleteSub(subId uuid.UUID) error {
	return service.subRepository.Delete(subId)
}

func (service SubscriptionService) GetSubsList(pageNumber int) ([]*response.SubResponse, error) {
	subList, err := service.subRepository.GetPage(pageNumber-1, 10)
	if err != nil {
		return make([]*response.SubResponse, 0), err
	}
	subResponses := make([]*response.SubResponse, len(subList), len(subList))
	for i, sub := range subList {
		subResponses[i] = response.NewSubResponse(sub)
	}
	return subResponses, nil
}

func (service SubscriptionService) GetSubsSum(filter *request.SubFilterRequest) (int, error) {
	startPeriod, err := parseDate(&filter.StartPeriod)
	if err != nil {
		return 0, fmt.Errorf("%w (start_date)", err)
	}
	endPeriod, err := parseDate(&filter.EndPeriod)
	if err != nil {
		return 0, fmt.Errorf("%w (end_date)", err)
	}

	submitedFiler := request.SubFilter{
		StartPeriod:  *startPeriod,
		EndPeriod:    *endPeriod,
		UserId:       filter.UserId,
		ServiceTitle: filter.ServiceTitle,
	}
	return service.subRepository.GetPriceSumByFilter(&submitedFiler)
}
