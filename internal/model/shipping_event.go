package model

type ShippingEvent struct {
	ID                  int64  `db:"id" json:"id"`
	Name                string `db:"name" json:"name"`
	ExpectedShipDate    int64  `db:"expected_ship_date" json:"expected_ship_date"`
	ActualShipDate      int64  `db:"actual_ship_date" json:"actual_ship_date"`
	ExpectedArrivalDate int64  `db:"expected_arrival_date" json:"expected_arrival_date"`
	ActualArrivalDate   int64  `db:"actual_arrival_date" json:"actual_arrival_date"`
	Status              string `db:"status" json:"status"`
	AuditField          string `db:"audit_field" json:"audit_field"`
}
