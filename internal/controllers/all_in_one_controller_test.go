package controllers

import (
	shippingevent "go-native-webserver/internal/service/shipping_event"
	"testing"
)

var mockController allInOneController

func TestMain(m *testing.M) {
	m.Run()
}

func TestPostSubscriptionSuccessResponse(t *testing.T) {
	mockController = allInOneController{
		shippingEventService: &shippingevent.MockShippingEventService{},
	}
	mockController.PostSubscription(nil, nil)
}
