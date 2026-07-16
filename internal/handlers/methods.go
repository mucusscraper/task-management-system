package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mucusscraper/task-management-system/internal/dto"
	"github.com/mucusscraper/task-management-system/internal/models"
	"github.com/mucusscraper/task-management-system/internal/service"
)

func handleError(ctx *gin.Context, err error) {
	if errors.Is(err, service.ErrUnauthorized) || errors.Is(err, service.ErrWorkerOnly) || errors.Is(err, service.ErrDifferentWorker) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, service.ErrTaskNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, service.ErrInvalidTransition) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

// CreateTask handles incoming requests to create a new task.
// Expects a JSON payload from a Supervisor.
func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	currentUser, _ := ctx.Get("currentUser")
	user, _ := currentUser.(models.User)
	var req dto.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	task, err := h.service.CreateTask(req, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, task)
}

// AssignTask handles request mapping for designating a worker to a task.
// Expects a Task ID in the URL and Worker ID in the body.
func (h *TaskHandler) AssignTask(ctx *gin.Context) {
	currentUser, _ := ctx.Get("currentUser")
	user, _ := currentUser.(models.User)

	var req dto.AssignTaskRequest
	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_task_id",
		})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err := h.service.AssignTask(taskID, req, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// UpdateTaskStatus processes Worker requests to change a task's state.
// Verifies status constraints before updating.
func (h *TaskHandler) UpdateTaskStatus(ctx *gin.Context) {
	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_task_id",
		})
		return
	}
	currentUser, _ := ctx.Get("currentUser")
	user, _ := currentUser.(models.User)
	var req dto.UpdateTaskStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	task, err := h.service.UpdateStatus(taskID, req, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// GetTasks retrieves a filtered task list based on the caller's context role.
func (h *TaskHandler) GetTasks(ctx *gin.Context) {
	currentUser, _ := ctx.Get("currentUser")
	user, _ := currentUser.(models.User)
	tasks, err := h.service.GetTasks(user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

// HealthCheck responds with the operational health of the application and its database link.
func (h *TaskHandler) HealthCheck(ctx *gin.Context) {
	if err := h.service.PingDatabase(); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
			"error":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}

// GetTaskByID fetches a single task profile from the store.
// Enforces role permissions on access.
func (h *TaskHandler) GetTaskByID(ctx *gin.Context) {
	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_task_id",
		})
		return
	}
	currentUser, _ := ctx.Get("currentUser")
	user, _ := currentUser.(models.User)
	task, err := h.service.GetTaskByID(taskID, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// GetNotifications returns a sorted log of notification messages for the calling Worker.
func (h *TaskHandler) GetNotifications(ctx *gin.Context) {
	currentUser, _ := ctx.Get("currentUser")
	user, _ := currentUser.(models.User)
	notifications, err := h.service.GetNotificationsByUser(user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, notifications)
}
