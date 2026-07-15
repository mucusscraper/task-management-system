package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mucusscraper/task-management-system/internal/handlers"
	"github.com/mucusscraper/task-management-system/internal/service"
)

func main() {
	server := gin.Default()
	tasks := server.Group("/tasks")
	taskService := service.NewTaskService()
	handler := handlers.NewTaskHandler(taskService)
	tasks.GET("", handler.GetTasks)
	tasks.POST("", handler.CreateTask)
	tasks.POST("/:id/assign", handler.AssignTask)
	tasks.PATCH("/:id/status", handler.UpdateTaskStatus)
	server.Run(":8080")
}
