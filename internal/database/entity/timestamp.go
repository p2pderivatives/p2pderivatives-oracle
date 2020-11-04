package entity

import "time"

// Timestamp entity to add db timestamp field
type Timestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
