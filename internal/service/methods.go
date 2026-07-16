package service

import (
	"github.com/mucusscraper/task-management-system/internal/dto"
	"github.com/mucusscraper/task-management-system/internal/models"
)

func (s *TaskService) CreateTask(req dto.CreateTaskRequest, user models.User) (*models.Task, error) {
	if string(user.Role) != string(models.Supervisor) {
		return nil, ErrUnauthorized
	}
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.Created,
	}
	query := `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(query, task.Title, task.Description, string(task.Status)).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) AssignTask(taskID int, req dto.AssignTaskRequest, user models.User) (*models.Task, error) {
	return nil, nil
}
