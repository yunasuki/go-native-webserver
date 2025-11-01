package shippingevent

import (
	"context"
	"errors"
	"fmt"
	"go-native-webserver/internal/apperror"
	"go-native-webserver/internal/dal"
	"go-native-webserver/internal/jobs"
	"go-native-webserver/internal/model"
	"go-native-webserver/internal/repositories"
	"go-native-webserver/internal/service"
)

type ShippingEventService interface {
	UpdateShippingEvent(ctx context.Context, eventID int64, status string) error
	AddUserToShippingEventSubscription(ctx context.Context, eventID int64, userID int64) error
}

type shippingEventService struct {
	service.BaseService
	userRepo                      repositories.UserRepository
	shippingEventRepo             repositories.ShippingEventRepository
	userShippingEventSubscription repositories.UserShippingEventSubscriptionRepository
}

// i want to use google's wire... native... urgh...
func NewShippingEventService() *shippingEventService {
	return &shippingEventService{
		BaseService:                   service.BaseService{},
		userRepo:                      repositories.NewUserRepository(),
		shippingEventRepo:             repositories.NewShippingEventRepository(),
		userShippingEventSubscription: repositories.NewUserShippingEventSubscriptionRepository(),
	}
}

func (s *shippingEventService) UpdateShippingEvent(ctx context.Context, eventID int64, status string) error {
	// Implement update logic here
	// 1. Update the shipping event
	expectedRecord := &model.ShippingEvent{
		ID:     eventID,
		Status: status,
	}
	err := s.shippingEventRepo.Update(ctx, expectedRecord)
	if err != nil {
		return apperror.APIError{
			Code:    500,
			Message: "Failed to update shipping event",
		}
	}

	// 2. Notify subscribed users (omitted for brevity)
	subscriptions, err := s.userShippingEventSubscription.ListByShippingEventID(eventID)
	if err != nil {
		if errors.Is(err, dal.ErrRecordNotFound) {
			return nil // no subscribers, nothing to do
		}
		return apperror.APIError{
			Code:    500,
			Message: "Internal Error - Failed to retrieve subscriptions",
		}
	}
	var userIDs []int64
	for _, sub := range subscriptions {
		userIDs = append(userIDs, sub.UserID)
	}
	newJob := &jobs.ShippingEventNotificationJob{
		EventID:   eventID,
		NewStatus: status,
		UserIDs:   userIDs,
	}

	// Push to database job records, let queue worker pick it up later
	fmt.Printf("Debug: Created job: %+v\n", newJob)

	return nil
}

func (s *shippingEventService) AddUserToShippingEventSubscription(ctx context.Context, eventID int64, userID int64) error {
	if eventID <= 0 || userID <= 0 {
		return apperror.APIError{
			Code:    404,
			Message: "Invalid eventID or userID",
		}
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, dal.ErrRecordNotFound) {
			return apperror.APIError{
				Code:    404,
				Message: "User invalid or user inactive",
			}
		}
		return apperror.APIError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	shippingEvent, err := s.shippingEventRepo.FindByID(ctx, eventID)
	if err != nil {
		if errors.Is(err, dal.ErrRecordNotFound) {
			return apperror.APIError{
				Code:    404,
				Message: "Shipping event not found",
			}
		}
		return apperror.APIError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	_, err = s.userShippingEventSubscription.CreateSubscription(user.ID, shippingEvent.ID)
	if err != nil {
		return apperror.APIError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return nil
}

type MockShippingEventService struct{}

func (s *MockShippingEventService) UpdateShippingEvent(ctx context.Context, eventID int64, status string) error {
	return nil
}

func (s *MockShippingEventService) AddUserToShippingEventSubscription(ctx context.Context, eventID int64, userID int64) error {
	// Mock implementation
	return nil
}
