package repositories

import (
	"context"
	"go-native-webserver/internal/dal"
	"go-native-webserver/internal/model"

	"gorm.io/gorm"
)

type ShippingEventRepository interface {
	FindByID(ctx context.Context, id int64) (*model.ShippingEvent, error)
	Create(ctx context.Context, event *model.ShippingEvent) error
	Update(ctx context.Context, event *model.ShippingEvent) error
	Delete(ctx context.Context, id int64) error
}

type shippingEventRepository struct {
	db *gorm.DB
}

func NewShippingEventRepository() ShippingEventRepository {
	return &shippingEventRepository{
		db: dal.DB, // assuming dal.DB is your gorm DB instance
	}
}

func (r *shippingEventRepository) FindByID(ctx context.Context, id int64) (*model.ShippingEvent, error) {
	var event model.ShippingEvent
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *shippingEventRepository) Create(ctx context.Context, event *model.ShippingEvent) error {
	// Implement the logic to create a new shipping event in the database
	return nil
}

func (r *shippingEventRepository) Update(ctx context.Context, event *model.ShippingEvent) error {
	if err := r.db.Where("id = ? ", event.ID).Updates(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *shippingEventRepository) Delete(ctx context.Context, id int64) error {
	// Implement the logic to delete a shipping event from the database
	return nil
}
