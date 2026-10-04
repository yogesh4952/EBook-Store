package models

import (
	"time"

	"gorm.io/gorm"
)

// Base mirrors gorm.Model field-for-field (same Go names, same GORM tags,
// same column names) so the DB schema is unchanged. The only difference is
// the json tags, which make the API emit snake_case like the rest of the app.
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
