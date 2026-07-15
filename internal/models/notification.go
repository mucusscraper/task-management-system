package models

import "time"

type Notification struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	UserID    int       `json:"user"`
	TaskID    int       `json:"task"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
