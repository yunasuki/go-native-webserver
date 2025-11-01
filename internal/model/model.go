package model

type AuditField struct {
	CreatedAt int64 `db:"created_at" json:"created_at"`
	CreatedBy int64 `db:"created_by" json:"created_by"`
	UpdatedAt int64 `db:"updated_at" json:"updated_at"`
	UpdatedBy int64 `db:"updated_by" json:"updated_by"`
	DeletedAt int64 `db:"deleted_at" json:"deleted_at"`
}
