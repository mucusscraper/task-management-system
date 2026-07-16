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

func TestAssignTask(t *testing.T) {
	db, err := database.NewPostgres()
	if err != nil {
		t.Skip("PostgreSQL not connecting")
		return
	}
	defer db.Close()

	service := NewTaskService(db)

	supervisor := models.User{ID: 1, Name: "Supervisor", Role: "SUPERVISOR"}
	worker := models.User{ID: 2, Name: "Worker 1", Role: "WORKER"}

	createReq := dto.CreateTaskRequest{
		Title:       "First Task",
		Description: "It will be used in tests for AssignTask",
	}
	createdTask, err := service.CreateTask(createReq, supervisor)
	if err != nil {
		t.Fatalf("Error creating task (CreateTask): %v", err)
	}

	t.Run("Should assign task successfully if user is Supervisor", func(t *testing.T) {
		assignReq := dto.AssignTaskRequest{
			WorkerID: worker.ID,
		}

		updatedTask, err := service.AssignTask(createdTask.ID, assignReq, supervisor)
		if err != nil {
			t.Fatalf("did not expect error, but got: %v", err)
		}

		if updatedTask == nil {
			t.Fatal("expected updated task, got nil")
		}

		if updatedTask.Status != "ASSIGNED" {
			t.Errorf("expected status ASSIGNED, got %q", updatedTask.Status)
		}

		if updatedTask.AssignedTo == nil || *updatedTask.AssignedTo != worker.ID {
			t.Errorf("expected task to be assigned to worker %d", worker.ID)
		}

		var notificationCount int
		err = db.QueryRow("SELECT COUNT(*) FROM notifications WHERE task_id = $1 AND user_id = $2", createdTask.ID, worker.ID).Scan(&notificationCount)
		if err != nil {
			t.Errorf("failed to verify notifications in database: %v", err)
		}
		if notificationCount == 0 {
			t.Error("expected a notification to be saved in database")
		}
	})

	t.Run("Should return ErrUnauthorized if assigning user is Worker", func(t *testing.T) {
		assignReq := dto.AssignTaskRequest{
			WorkerID: worker.ID,
		}

		_, err := service.AssignTask(createdTask.ID, assignReq, worker)
		if err == nil {
			t.Fatal("expected authorization error, but did not get one")
		}

		if !errors.Is(err, ErrUnauthorized) {
			t.Errorf("expected error %v, got %v", ErrUnauthorized, err)
		}
	})

	t.Run("Should return ErrTaskNotFound if task ID does not exist", func(t *testing.T) {
		assignReq := dto.AssignTaskRequest{
			WorkerID: worker.ID,
		}

		_, err := service.AssignTask(99999, assignReq, supervisor) // Non-existent ID
		if err == nil {
			t.Fatal("expected task not found error, but did not get one")
		}

		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
		}
	})

	t.Run("Should return ErrInvalidTransition if task is already assigned or in progress", func(t *testing.T) {
		assignReq := dto.AssignTaskRequest{
			WorkerID: worker.ID,
		}

		// Attempting to assign the same task again (already 'ASSIGNED' from the first test case)
		_, err := service.AssignTask(createdTask.ID, assignReq, supervisor)
		if err == nil {
			t.Fatal("expected invalid transition error, but did not get one")
		}

		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected error %v, got %v", ErrInvalidTransition, err)
		}
	})
}

func TestUpdateStatus(t *testing.T) {
	db, err := database.NewPostgres()
	if err != nil {
		t.Skip("PostgreSQL not connecting")
		return
	}
	defer db.Close()

	service := NewTaskService(db)

	supervisor := models.User{ID: 1, Name: "Supervisor", Role: "SUPERVISOR"}
	worker1 := models.User{ID: 2, Name: "Worker 1", Role: "WORKER"}
	worker2 := models.User{ID: 3, Name: "Worker 2", Role: "WORKER"}

	createReq := dto.CreateTaskRequest{
		Title:       "Task to be processed",
		Description: "Testing status flow",
	}
	task, err := service.CreateTask(createReq, supervisor)
	if err != nil {
		t.Fatalf("failed to setup task: %v", err)
	}

	_, err = service.AssignTask(task.ID, dto.AssignTaskRequest{WorkerID: worker1.ID}, supervisor)
	if err != nil {
		t.Fatalf("failed to assign task: %v", err)
	}

	t.Run("Should return ErrWorkerOnly if supervisor tries to update status", func(t *testing.T) {
		_, err := service.UpdateStatus(task.ID, dto.UpdateTaskStatusRequest{}, supervisor)
		if !errors.Is(err, ErrWorkerOnly) {
			t.Errorf("expected %v, got %v", ErrWorkerOnly, err)
		}
	})

	t.Run("Should return ErrUnauthorized if a different worker tries to update status", func(t *testing.T) {
		_, err := service.UpdateStatus(task.ID, dto.UpdateTaskStatusRequest{}, worker2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Errorf("expected %v, got %v", ErrUnauthorized, err)
		}
	})

	t.Run("Should transition successfully from ASSIGNED to IN_PROGRESS", func(t *testing.T) {
		updated, err := service.UpdateStatus(task.ID, dto.UpdateTaskStatusRequest{}, worker1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Status != "IN_PROGRESS" {
			t.Errorf("expected status IN_PROGRESS, got %s", updated.Status)
		}
	})

	t.Run("Should transition successfully from IN_PROGRESS to COMPLETED", func(t *testing.T) {
		updated, err := service.UpdateStatus(task.ID, dto.UpdateTaskStatusRequest{}, worker1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Status != "COMPLETED" {
			t.Errorf("expected status COMPLETED, got %s", updated.Status)
		}
	})

	t.Run("Should return ErrInvalidTransition if task is already COMPLETED", func(t *testing.T) {
		_, err := service.UpdateStatus(task.ID, dto.UpdateTaskStatusRequest{}, worker1)
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected %v, got %v", ErrInvalidTransition, err)
		}
	})
}

func TestGetTasks(t *testing.T) {
	db, err := database.NewPostgres()
	if err != nil {
		t.Skip("PostgreSQL not connecting")
		return
	}
	defer db.Close()

	service := NewTaskService(db)

	supervisor := models.User{ID: 1, Name: "Supervisor", Role: "SUPERVISOR"}
	worker1 := models.User{ID: 2, Name: "Worker 1", Role: "WORKER"}
	worker2 := models.User{ID: 3, Name: "Worker 2", Role: "WORKER"}

	_, err = db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("failed to clean tasks table: %v", err)
	}

	t1, err := service.CreateTask(dto.CreateTaskRequest{Title: "Task One"}, supervisor)
	if err != nil {
		t.Fatalf("failed to setup task 1: %v", err)
	}
	_, err = service.CreateTask(dto.CreateTaskRequest{Title: "Task Two"}, supervisor)
	if err != nil {
		t.Fatalf("failed to setup task 2: %v", err)
	}

	// Assign Task One to Worker 1 (Task Two remains unassigned/created)
	_, err = service.AssignTask(t1.ID, dto.AssignTaskRequest{WorkerID: worker1.ID}, supervisor)
	if err != nil {
		t.Fatalf("failed to assign task 1: %v", err)
	}

	t.Run("Should list all tasks in the system if user is Supervisor", func(t *testing.T) {
		tasks, err := service.GetTasks(supervisor)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tasks) != 2 {
			t.Errorf("expected exactly 2 tasks for supervisor, got %d", len(tasks))
		}
	})

	t.Run("Should only list tasks assigned to the specific Worker", func(t *testing.T) {
		tasks, err := service.GetTasks(worker1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tasks) != 1 {
			t.Errorf("expected exactly 1 task for worker 1, got %d", len(tasks))
		}

		if tasks[0].ID != t1.ID {
			t.Errorf("expected task ID %d, got %d", t1.ID, tasks[0].ID)
		}
	})

	t.Run("Should return an empty slice if Worker has no tasks assigned", func(t *testing.T) {
		tasks, err := service.GetTasks(worker2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(tasks) != 0 {
			t.Errorf("expected 0 tasks for worker 2, got %d", len(tasks))
		}
	})
}

func TestGetTaskByID(t *testing.T) {
	db, err := database.NewPostgres()
	if err != nil {
		t.Skip("PostgreSQL not connecting")
		return
	}
	defer db.Close()

	service := NewTaskService(db)

	supervisor := models.User{ID: 1, Name: "Supervisor", Role: "SUPERVISOR"}
	worker1 := models.User{ID: 2, Name: "Worker 1", Role: "WORKER"}
	worker2 := models.User{ID: 3, Name: "Worker 2", Role: "WORKER"}

	task, err := service.CreateTask(dto.CreateTaskRequest{Title: "Single Task test"}, supervisor)
	if err != nil {
		t.Fatalf("failed to setup task: %v", err)
	}

	_, err = service.AssignTask(task.ID, dto.AssignTaskRequest{WorkerID: worker1.ID}, supervisor)
	if err != nil {
		t.Fatalf("failed to assign task: %v", err)
	}

	t.Run("Supervisor should be able to fetch any task by ID", func(t *testing.T) {
		fetched, err := service.GetTaskByID(task.ID, supervisor)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fetched.ID != task.ID {
			t.Errorf("expected task ID %d, got %d", task.ID, fetched.ID)
		}
	})

	t.Run("Assigned Worker should be able to fetch their own task", func(t *testing.T) {
		fetched, err := service.GetTaskByID(task.ID, worker1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fetched.ID != task.ID {
			t.Errorf("expected task ID %d, got %d", task.ID, fetched.ID)
		}
	})

	t.Run("Should return ErrUnauthorized if worker is not assigned to the task", func(t *testing.T) {
		_, err := service.GetTaskByID(task.ID, worker2)
		if !errors.Is(err, ErrUnauthorized) {
			t.Errorf("expected error %v, got %v", ErrUnauthorized, err)
		}
	})

	t.Run("Should return ErrTaskNotFound for invalid task IDs", func(t *testing.T) {
		_, err := service.GetTaskByID(99999, supervisor)
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
		}
	})
}

func TestGetNotificationsByUser(t *testing.T) {
	db, err := database.NewPostgres()
	if err != nil {
		t.Skip("PostgreSQL not connecting")
		return
	}
	defer db.Close()

	service := NewTaskService(db)

	supervisor := models.User{ID: 1, Name: "Supervisor", Role: "SUPERVISOR"}
	worker1 := models.User{ID: 2, Name: "Worker 1", Role: "WORKER"}

	_, _ = db.Exec("DELETE FROM notifications")
	_, _ = db.Exec("DELETE FROM tasks")

	task, err := service.CreateTask(dto.CreateTaskRequest{Title: "Notify Test"}, supervisor)
	if err != nil {
		t.Fatalf("failed to setup task: %v", err)
	}

	_, err = service.AssignTask(task.ID, dto.AssignTaskRequest{WorkerID: worker1.ID}, supervisor)
	if err != nil {
		t.Fatalf("failed to assign task: %v", err)
	}

	t.Run("Worker should retrieve their notification history", func(t *testing.T) {
		notifications, err := service.GetNotificationsByUser(worker1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(notifications) != 1 {
			t.Errorf("expected 1 notification, got %d", len(notifications))
		}
		if notifications[0].UserID != worker1.ID {
			t.Errorf("expected notification for user %d, got user %d", worker1.ID, notifications[0].UserID)
		}
	})

	t.Run("Should return ErrWorkerOnly if Supervisor tries to get notifications", func(t *testing.T) {
		_, err := service.GetNotificationsByUser(supervisor)
		if !errors.Is(err, ErrWorkerOnly) {
			t.Errorf("expected error %v, got %v", ErrWorkerOnly, err)
		}
	})
}
