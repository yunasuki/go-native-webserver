package repositories

import (
	"go-native-webserver/internal/dal"
	"go-native-webserver/internal/model"
)

type UserShippingEventSubscriptionRepository interface {
	// Define methods for managing user shipping event subscriptions
	CreateSubscription(userID int64, shippingEventID int64) (*model.UserShippingEventSubscription, error)
	ListByShippingEventID(eventID int64) ([]model.UserShippingEventSubscription, error)
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

func (repo *userShippingEventSubscriptionRepo) ListByShippingEventID(eventID int64) ([]model.UserShippingEventSubscription, error) {
	var subscriptions []model.UserShippingEventSubscription
	err := repo.db.Model(&model.UserShippingEventSubscription{}).Where("shipping_event_id = ?", eventID).Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}
