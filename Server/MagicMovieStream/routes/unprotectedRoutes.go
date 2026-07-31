package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetupUnprotectedRoutes are the routes anyone can hit without a token.
func SetupUnprotectedRoutes(router *gin.RouterGroup, client *mongo.Client) {
	router.GET("/movies", controller.GetMovies(client))
	router.GET("/movie/:imdb_id", controller.GetMovie(client))
	router.GET("/genres", controller.GetGenres(client))
	router.GET("/media", controller.ListMedia())
	router.GET("/media/*filepath", controller.StreamMedia())
}

// SetupAuthRoutes are the routes that create or rotate sessions.
// They get a stricter rate limit than the public read routes.
func SetupAuthRoutes(router *gin.RouterGroup, client *mongo.Client) {
	router.POST("/user", controller.RegisterUser(client))
	router.POST("/login", controller.LoginUser(client))
	router.POST("/refreshtoken", controller.RefreshTokenHandler(client))
}
