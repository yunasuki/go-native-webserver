package jobs

import (
	"context"
	"fmt"
	"go-native-webserver/pkg/queue"
)

var _ queue.Job = (*ShippingEventNotificationJob)(nil)

type ShippingEventNotificationJob struct {
	// Define fields relevant to the shipping event notification job
	EventID   int64
	NewStatus string
	UserIDs   []int64 // list of user IDs to notify
}

func (job *ShippingEventNotificationJob) Process(ctx context.Context) error {
	fmt.Printf("Processing ShippingEventNotificationJob for EventID: %d, NewStatus: %s\n", job.EventID, job.NewStatus)
	for _, userID := range job.UserIDs {
		fmt.Printf("Send Email/App notification To UserID: %d about EventID: %d status change to %s\n", userID, job.EventID, job.NewStatus)
	}
	return nil
}
