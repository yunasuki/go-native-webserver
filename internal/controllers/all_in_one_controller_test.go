package controllers

import (
	"fmt"
	shippingevent "go-native-webserver/internal/service/shipping_event"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mockController allInOneController

func TestMain(m *testing.M) {
	m.Run()
}

// hmm... native testing....
func TestPostSubscription(t *testing.T) {
	mockController = allInOneController{
		shippingEventService: &shippingevent.MockShippingEventService{},
	}

	testData := []struct {
		UserID         int64 // maybe use string to make more stupid request body tests....
		EventID        int64
		ExpectedStatus int
	}{
		{UserID: 1, EventID: 100, ExpectedStatus: http.StatusAccepted},
		{UserID: 1, EventID: 100, ExpectedStatus: http.StatusBadRequest}, // depend if allow same subscription more then once...
		{UserID: 2, EventID: 200, ExpectedStatus: http.StatusAccepted},
		{UserID: 0, EventID: 100, ExpectedStatus: http.StatusBadRequest},
		{UserID: 1, EventID: 0, ExpectedStatus: http.StatusBadRequest},
	}

	for _, td := range testData {
		jsonBody := fmt.Sprintf(`{"user_id": %d, "event_id": %d}`, td.UserID, td.EventID)
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		t.Run(fmt.Sprintf("UserID:%d_EventID:%d", td.UserID, td.EventID), func(t *testing.T) {
			mockController.PostSubscription(w, req)
			if w.Result().StatusCode != td.ExpectedStatus {
				t.Errorf("Expected status %d, got %d", td.ExpectedStatus, w.Result().StatusCode)
			}
		})
	}
}
