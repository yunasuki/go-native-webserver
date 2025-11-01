package shippingevent

import (
	"context"
	"errors"
	"go-native-webserver/internal/apperror"
	"go-native-webserver/internal/dal"
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

func (m *MockShippingEventService) UpdateShippingEvent(ctx context.Context, eventID int64, status string) error {
	return nil
}

func (m *MockShippingEventService) AddUserToShippingEventSubscription(ctx context.Context, eventID int64, userID int64) error {
	// Mock implementation
	return nil
}
