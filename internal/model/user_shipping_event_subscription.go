package model

// if record exist, when ShippingEvent updated, notify user
type UserShippingEventSubscription struct {
	ID              int64  `db:"id" json:"id"`
	UserID          int64  `db:"user_id" json:"user_id"`
	ShippingEventID int64  `db:"shipping_event_id" json:"shipping_event_id"`
	AuditField      string `db:"audit_field" json:"audit_field"`
}
