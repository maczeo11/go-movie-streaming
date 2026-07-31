package controllers

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var serverStartedAt = time.Now()

func Health(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 2*time.Second)
		defer cancel()

		dbUp := client.Ping(ctx, nil) == nil

		status := http.StatusOK
		statusText := "ok"
		if !dbUp {
			status = http.StatusServiceUnavailable
			statusText = "degraded"
		}

		c.JSON(status, gin.H{
			"status":         statusText,
			"uptime_seconds": int(time.Since(serverStartedAt).Seconds()),
			"database": gin.H{
				"connected": dbUp,
			},
			"go_version": runtime.Version(),
			"env":        gin.Mode(),
		})
	}
}
