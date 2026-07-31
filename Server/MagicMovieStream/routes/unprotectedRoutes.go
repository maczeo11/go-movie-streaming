package routes

import (
	controller "github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnprotectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.GET("/movies", controller.GetMovies(client))
	router.GET("/genres", controller.GetGenres(client))
	router.POST("/user", controller.RegisterUser(client))
	router.POST("/login", controller.LoginUser(client))
	router.POST("/logout", controller.LogoutHandler(client))
	router.POST("/refreshtoken", controller.RefreshTokenHandler(client))
}