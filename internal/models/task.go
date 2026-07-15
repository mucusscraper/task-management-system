package models

type TaskStatus string

const (
	Created    TaskStatus = "CREATED"
	Assigned   TaskStatus = "ASSIGNED"
	InProgress TaskStatus = "IN_PROGRESS"
	Completed  TaskStatus = "COMPLETED"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	AssignedTo  *int       `json:"assigned_to"`
}
