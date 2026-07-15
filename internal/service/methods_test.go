package service

import (
	"testing"

	"github.com/mucusscraper/task-management-system/internal/dto"
	"github.com/mucusscraper/task-management-system/internal/models"
)

func TestCreateTask(t *testing.T) {
	service := NewTaskService()
	req := dto.CreateTaskRequest{
		Title:       "Write complete documentation",
		Description: "Complete the README",
	}
	task, err := service.CreateTask(req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.Title != req.Title {
		t.Errorf("expected title %q, got %q", req.Title, task.Title)
	}
	if task.Description != req.Description {
		t.Errorf("expected description %q, got %q", req.Description, task.Description)
	}
	if task.Status != models.Created {
		t.Errorf("expected status %q, got %q", models.Created, task.Status)
	}
}
