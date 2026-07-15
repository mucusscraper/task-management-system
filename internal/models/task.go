package models

import "time"

type TaskStatus string

const (
	Created    TaskStatus = "CREATED"
	Assigned   TaskStatus = "ASSIGNED"
	InProgress TaskStatus = "IN_PROGRESS"
	Completed  TaskStatus = "COMPLETED"
)

type Role string

const (
	Supervisor Role = "SUPERVISOR"
	Worker     Role = "WORKER"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	AssignedTo  *int       `json:"assigned_to"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Notification struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	UserID    int       `json:"user"`
	TaskID    int       `json:"task"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
