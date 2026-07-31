package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/controllers"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/middleware"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupRouter(client *mongo.Client) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(slog.Default()))
	router.Use(cors.New(corsConfig()))

	router.GET("/health", controllers.Health(client))

	// Public reads - generous limit.
	public := router.Group("/")
	public.Use(middleware.RateLimit(newRateLimiter("API_RATE", "API_BURST", 2, 120)))
	routes.SetupUnprotectedRoutes(public, client)

	// Auth endpoints - tight limit to slow down password guessing.
	auth := router.Group("/")
	auth.Use(middleware.RateLimit(newRateLimiter("AUTH_RATE", "AUTH_BURST", 0.1, 10)))
	routes.SetupAuthRoutes(auth, client)

	// Everything below needs a valid token.
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	routes.SetupProtectedRoutes(protected, client)

	admin := protected.Group("/")
	admin.Use(middleware.AdminOnly())
	routes.SetupAdminRoutes(admin, client)

	return router
}

func corsConfig() cors.Config {
	config := cors.DefaultConfig()

	allowed := os.Getenv("ALLOWED_ORIGINS")
	if allowed != "" {
		var origins []string
		for _, o := range strings.Split(allowed, ",") {
			if o = strings.TrimSpace(o); o != "" {
				origins = append(origins, o)
			}
		}
		config.AllowOrigins = origins
	} else {
		config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:4173"}
	}

	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length", "X-Request-ID"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	return config
}

func newRateLimiter(rateKey, burstKey string, defaultRate, defaultBurst float64) *middleware.RateLimiter {
	return middleware.NewRateLimiter(
		envFloat(rateKey, defaultRate),
		envFloat(burstKey, defaultBurst),
		time.Minute,
	)
}

func envFloat(key string, fallback float64) float64 {
	if raw := os.Getenv(key); raw != "" {
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			return v
		}
	}
	return fallback
}
