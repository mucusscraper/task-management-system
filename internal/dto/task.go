package dto

import "github.com/mucusscraper/task-management-system/internal/models"

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type AssignTaskRequest struct {
	WorkerID int `json:"worker_id" binding:"required"`
}

type UpdateTaskStatusRequest struct {
	Status models.TaskStatus `json:"status" binding:"required"`
}
