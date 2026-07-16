package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
		INSERT INTO tasks (title, description, status, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(query, task.Title, task.Description, string(task.Status), user.ID).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetTasks(user models.User) ([]models.Task, error) {
	var rows *sql.Rows
	var err error
	if string(user.Role) == string(models.Worker) {
		getQuery := `SELECT id, title, description, status, assigned_to, created_at, updated_at 
			FROM tasks 
			WHERE assigned_to=$1
		`
		rows, err = s.db.Query(getQuery, user.ID)
	} else if string(user.Role) == string(models.Supervisor) {
		getQuery := `SELECT id, title, description, status, assigned_to, created_at, updated_at 
			FROM tasks 
		`
		rows, err = s.db.Query(getQuery)
	} else {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.AssignedTo,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) UpdateStatus(taskID int, req dto.UpdateTaskStatusRequest, user models.User) (*models.Task, error) {
	if string(user.Role) != string(models.Worker) {
		return nil, ErrWorkerOnly
	}
	var currentStatus string
	var assignedTo *int
	checkQuery := `SELECT status, assigned_to FROM tasks WHERE id=$1`
	err := s.db.QueryRow(checkQuery, taskID).Scan(&currentStatus, &assignedTo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	if assignedTo == nil || *assignedTo != user.ID {
		return nil, ErrDifferentWorker
	}
	if currentStatus != string(models.Assigned) && currentStatus != string(models.InProgress) {
		return nil, ErrInvalidTransition
	}
	task := &models.Task{}
	if currentStatus == string(models.Assigned) {
		updateQuery := `
			UPDATE tasks
				SET status = 'IN_PROGRESS',
				updated_at = NOW()
			WHERE id=$1
			RETURNING id,title,description, status,assigned_to,created_at,updated_at
		`
		err = s.db.QueryRow(updateQuery, taskID).Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssignedTo,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notificationQuery := `
			INSERT INTO notifications (task_id,user_id,message)
			VALUES ($1,$2,$3)
		`
		notificationMessage := fmt.Sprintf("Task %d of user %d now in progress", taskID, user.ID)
		_, err = s.db.Exec(notificationQuery, taskID, user.ID, notificationMessage)
		if err != nil {
			log.Printf("failed to save notification: %v", err)
		}
		fmt.Println(notificationMessage)
	} else if currentStatus == string(models.InProgress) {
		updateQuery := `
			UPDATE tasks
				SET status = 'COMPLETED',
				updated_at = NOW()
			WHERE id=$1
			RETURNING id,title,description, status,assigned_to,created_at,updated_at
		`
		err = s.db.QueryRow(updateQuery, taskID).Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssignedTo,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notificationQuery := `
			INSERT INTO notifications (task_id,user_id,message)
			VALUES ($1,$2,$3)
		`
		notificationMessage := fmt.Sprintf("Task %d of user %d now complete", taskID, user.ID)
		_, err = s.db.Exec(notificationQuery, taskID, user.ID, notificationMessage)
		if err != nil {
			log.Printf("failed to save notification: %v", err)
		}
		fmt.Println(notificationMessage)
	}
	return task, nil
}

func (s *TaskService) AssignTask(taskID int, req dto.AssignTaskRequest, user models.User) (*models.Task, error) {
	if string(user.Role) != string(models.Supervisor) {
		return nil, ErrUnauthorized
	}
	var currentStatus string
	checkQuery := `SELECT status FROM tasks WHERE id=$1`
	err := s.db.QueryRow(checkQuery, taskID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	if currentStatus != string(models.Created) {
		return nil, ErrInvalidTransition
	}
	task := &models.Task{}
	updateQuery := `
		UPDATE tasks
			SET status = 'ASSIGNED',
			assigned_to = $1,
			updated_at = NOW()
		WHERE id=$2
		RETURNING id,title,description, status,assigned_to,created_at,updated_at
	`
	err = s.db.QueryRow(updateQuery, req.WorkerID, taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssignedTo,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	notificationQuery := `
		INSERT INTO notifications (task_id,user_id,message)
		VALUES ($1,$2,$3)
	`
	notificationMessage := fmt.Sprintf("Task %d assigned to worker %d", taskID, req.WorkerID)
	_, err = s.db.Exec(notificationQuery, taskID, req.WorkerID, notificationMessage)
	if err != nil {
		log.Printf("failed to save notification: %v", err)
	}
	fmt.Println(notificationMessage)
	return task, nil
}
