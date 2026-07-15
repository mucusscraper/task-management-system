package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		headerValue    string
		expectedStatus int
	}{
		{
			name:           "Must fail if header X-User-Id is missing",
			headerValue:    "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Must fail if header X-User-Id is non-numeric",
			headerValue:    "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			r.Use(AuthMiddleware(nil))

			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-User-Id", tt.headerValue)
			}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("esperava status %d, mas obteve %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
