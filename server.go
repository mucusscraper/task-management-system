package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	tasks := server.Group("/tasks")
	tasks.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Supervisor: all tasks; Worker: own tasks only",
		})
	})
	tasks.POST("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Supervisor — create task (status: CREATED)",
		})
	})
	tasks.POST("/:id/assign", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{
			"task_id": id,
			"message": "Supervisor — assign to worker (status: ASSIGNED)",
		})
	})
	tasks.PATCH("/:id/status", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{
			"task_id": id,
			"message": "Worker — set status to IN_PROGRESS or COMPLETED",
		})
	})
	server.Run(":8080")
}
