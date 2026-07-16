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
