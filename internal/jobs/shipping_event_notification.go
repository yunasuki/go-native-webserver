package jobs

import (
	"context"
	"go-native-webserver/pkg/queue"
	"time"
)

var _ queue.Job = (*ShippingEventNotificationJob)(nil)

type ShippingEventNotificationJob struct {
	// Define fields relevant to the shipping event notification job
	OrderID      string
	CustomerID   string
	ShippingDate time.Time
	Status       string
}

func (job *ShippingEventNotificationJob) Process(ctx context.Context) error {

	return nil
}
