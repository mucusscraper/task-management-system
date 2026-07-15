package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mucusscraper/task-management-system/internal/handlers"
)

func main() {
	server := gin.Default()
	tasks := server.Group("/tasks")
	handler := handlers.NewTaskHandler()
	tasks.GET("", handler.GetTasks)
	tasks.POST("", handler.CreateTask)
	tasks.POST("/:id/assign", handler.AssignTask)
	tasks.PATCH("/:id/status", handler.UpdateTaskStatus)
	server.Run(":8080")
}
