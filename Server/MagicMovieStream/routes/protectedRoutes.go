package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetupProtectedRoutes require a valid access token.
func SetupProtectedRoutes(router *gin.RouterGroup, client *mongo.Client) {
	router.POST("/logout", controller.LogoutHandler(client))
	router.GET("/profile", controller.GetProfile(client))
	router.PATCH("/profile/genres", controller.UpdateProfileGenres(client))
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
}

// SetupAdminRoutes are protected AND restricted to the ADMIN role.
func SetupAdminRoutes(router *gin.RouterGroup, client *mongo.Client) {
	router.POST("/addmovie", controller.AddMovie(client))
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))
}
