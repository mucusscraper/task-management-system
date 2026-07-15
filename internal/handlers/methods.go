package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mucusscraper/task-management-system/internal/dto"
)

func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	var req dto.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	task, err := h.service.CreateTask(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) AssignTask(ctx *gin.Context) {
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
	task, err := h.service.AssignTask(taskID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTaskStatus(ctx *gin.Context) {
	var req dto.UpdateTaskStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": req.Status,
	})
}

func (h *TaskHandler) GetTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "list tasks",
	})
}
