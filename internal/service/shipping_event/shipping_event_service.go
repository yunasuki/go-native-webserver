package shippingevent

import (
	"context"
	"go-native-webserver/internal/service"
)

type ShippingEventService interface {
	UpdateShippingEvent(ctx context.Context, eventID int64, status string) error
	AddUserToShippingEventSubscription(ctx context.Context, eventID int64, userID int64) error
}

type shippingEventService struct {
	service.BaseService
}

// i want to use google's wire... native... urgh...
func NewShippingEventService() *shippingEventService {
	return &shippingEventService{
		BaseService: service.BaseService{},
	}
}

func (s *shippingEventService) UpdateShippingEvent(ctx context.Context, eventID int64, status string) error {
	// Implement update logic here
	return nil
}

func (s *shippingEventService) AddUserToShippingEventSubscription(ctx context.Context, eventID int64, userID int64) error {
	// Implement subscription logic here
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
