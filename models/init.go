package models

import "time"

type ModelMixin struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
