package service

import (
	"errors"
	"testing"

	"github.com/mucusscraper/task-management-system/internal/database"
	"github.com/mucusscraper/task-management-system/internal/dto"
	"github.com/mucusscraper/task-management-system/internal/models"
)

func TestCreateTask(t *testing.T) {
	db, err := database.NewPostgres()
	if err != nil {
		t.Skip("Pulando teste: PostgreSQL local não está rodando ou configurado")
		return
	}
	defer db.Close()

	service := NewTaskService(db)
	supervisor := models.User{ID: 1, Name: "Supervisor", Role: "SUPERVISOR"}
	worker := models.User{ID: 2, Name: "Worker", Role: "WORKER"}
	req := dto.CreateTaskRequest{
		Title:       "Testar persistência no banco",
		Description: "Garantir que a query de INSERT está funcionando",
	}
	t.Run("Must create task with success if supervisor role", func(t *testing.T) {
		task, err := service.CreateTask(req, supervisor)
		if err != nil {
			t.Fatalf("didn't expect error, but got: %v", err)
		}
		if task.ID == 0 {
			t.Error("expected a real ID, got 0")
		}
		if task.Title != req.Title {
			t.Errorf("expected title %q, got %q", req.Title, task.Title)
		}
		if task.Status != models.Created {
			t.Errorf("expected initial status CREATED, got %q", task.Status)
		}
	})
	t.Run("Must return ErrUnauthorized if worker role", func(t *testing.T) {
		task, err := service.CreateTask(req, worker)

		if err == nil {
			t.Fatal("expected auth error, but didn't got")
		}
		if !errors.Is(err, ErrUnauthorized) {
			t.Errorf("expected error %q, but got %q", ErrUnauthorized, err)
		}
		if task != nil {
			t.Error("didn't expect any task to be created")
		}
	})
}
