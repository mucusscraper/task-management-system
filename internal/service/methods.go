package service

import (
	"github.com/mucusscraper/task-management-system/internal/dto"
	"github.com/mucusscraper/task-management-system/internal/models"
)

func (s *TaskService) CreateTask(req dto.CreateTaskRequest) (*models.Task, error) {

	task := &models.Task{
		ID:          1,
		Title:       req.Title,
		Description: req.Description,
		Status:      models.Created,
	}

	return task, nil
}

func (s *TaskService) AssignTask(taskID int, req dto.AssignTaskRequest) (*models.Task, error) {
	return nil, nil
}
