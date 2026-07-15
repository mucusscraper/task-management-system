package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mucusscraper/task-management-system/internal/models"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.GetHeader("X-User-Id")
		if userIDStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "X-User-Id header is required",
			})
			ctx.Abort()
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user id format",
			})
			ctx.Abort()
			return
		}
		var user models.User
		query := `SELECT id, name, role FROM users WHERE id=$1`
		err = db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Role)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found in seed data",
			})
			ctx.Abort()
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			ctx.Abort()
			return
		}
		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
