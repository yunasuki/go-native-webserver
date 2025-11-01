package controllers

import (
	"context"
	shippingevent "go-native-webserver/internal/service/shipping_event"
	"net/http"
)

type AllInOneController interface {
	PostSubscription(w http.ResponseWriter, r *http.Request)
	GetPublicHoliday(w http.ResponseWriter, r *http.Request)
}

type allInOneController struct {
	shippingEventService shippingevent.ShippingEventService
}

func NewAllInOneController() AllInOneController {
	return &allInOneController{
		shippingEventService: shippingevent.NewShippingEventService(),
	}
}

func (c *allInOneController) PostLogin(w http.ResponseWriter, r *http.Request) {

}

func (c *allInOneController) PostSubscription(w http.ResponseWriter, r *http.Request) {
	err := c.shippingEventService.AddUserToShippingEventSubscription(context.Background(), 1, 1)
	if err != nil {
		ResponseError(w, err)
		return
	}
	ResponseSuccessJSON(w, http.StatusAccepted, map[string]string{"status": "subscribed"})
}

func (c *allInOneController) GetPublicHoliday(w http.ResponseWriter, r *http.Request) {
	// Implement the logic for fetching public holidays.
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("GetPublicHoliday not implemented"))
}
