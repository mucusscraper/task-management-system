package models

type Role string

const (
	Supervisor Role = "SUPERVISOR"
	Worker     Role = "WORKER"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}
