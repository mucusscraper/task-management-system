package handlers

import (
	"net/http"

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
	ctx.JSON(http.StatusCreated, gin.H{
		"title": req.Title,
	})
}

func (h *TaskHandler) AssignTask(ctx *gin.Context) {
	var req dto.AssignTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"worker_id": req.WorkerID,
	})
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

}
