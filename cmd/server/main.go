package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mucusscraper/task-management-system/internal/database"
	"github.com/mucusscraper/task-management-system/internal/handlers"
	"github.com/mucusscraper/task-management-system/internal/service"
)

func main() {
	db, err := database.NewPostgres("migrations")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	server := gin.Default()
	tasks := server.Group("/tasks")
	taskService := service.NewTaskService(db)
	handler := handlers.NewTaskHandler(taskService)
	server.GET("/health", handler.HealthCheck)
	tasks.Use(handlers.AuthMiddleware(db))
	{
		tasks.GET("", handler.GetTasks)
		tasks.POST("", handler.CreateTask)
		tasks.POST("/:id/assign", handler.AssignTask)
		tasks.PATCH("/:id/status", handler.UpdateTaskStatus)
		tasks.GET("/notifications", handler.GetNotifications)
		tasks.GET("/:id", handler.GetTaskByID)
	}
	server.Run(":8080")
}
