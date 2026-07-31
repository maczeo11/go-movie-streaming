package routes

import (
	controller "github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetupProtectedRoutes require a valid access token.
func SetupProtectedRoutes(router *gin.RouterGroup, client *mongo.Client) {
	router.POST("/logout", controller.LogoutHandler(client))
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
}

// SetupAdminRoutes are protected AND restricted to the ADMIN role.
func SetupAdminRoutes(router *gin.RouterGroup, client *mongo.Client) {
	router.POST("/addmovie", controller.AddMovie(client))
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))
}
