package repositories

import (
	"go-native-webserver/internal/dal"
	"go-native-webserver/internal/model"
)

type UserShippingEventSubscriptionRepository interface {
	// Define methods for managing user shipping event subscriptions
	CreateSubscription(userID int64, shippingEventID int64) (*model.UserShippingEventSubscription, error)
	DeleteSubscription(subscriptionID int64) error
	GetSubscriptionsByEventID(eventID int64) ([]int64, error) // returns list of UserIDs
}

type userShippingEventSubscriptionRepo struct {
	db dal.DatabaseConnection
}

func NewUserShippingEventSubscriptionRepository() UserShippingEventSubscriptionRepository {
	return &userShippingEventSubscriptionRepo{
		db: dal.GetDB(),
	}
}

func (repo *userShippingEventSubscriptionRepo) CreateSubscription(userID int64, shippingEventID int64) (*model.UserShippingEventSubscription, error) {
	newRecord := &model.UserShippingEventSubscription{
		UserID:          userID,
		ShippingEventID: shippingEventID,
	}
	repo.db.Model(&model.UserShippingEventSubscription{}).Where("user_id = ? AND shipping_event_id = ? ", userID, shippingEventID).FirstOrCreate(newRecord)
	return newRecord, nil
}

func (repo *userShippingEventSubscriptionRepo) DeleteSubscription(subscriptionID int64) error {
	// Implement the logic to delete a subscription
	return nil
}

func (repo *userShippingEventSubscriptionRepo) GetSubscriptionsByEventID(eventID int64) ([]int64, error) {
	// Implement the logic to get user IDs subscribed to a specific shipping event
	return []int64{}, nil
}
