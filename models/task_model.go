package models

import "time"

// user task struct
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   string `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      int64 `json:"userId"`
}
