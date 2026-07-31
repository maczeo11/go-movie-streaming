package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/database"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/model"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetProfile(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var user model.User
		err = database.OpenCollection("users", client).FindOne(ctx, bson.D{{Key: "user_id", Value: userId}}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, model.UserResponse{
			UserId:          user.UserID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Role:            user.Role,
			FavouriteGenres: user.FavouriteGenres,
		})
	}
}

func UpdateProfileGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req struct {
			FavouriteGenres []model.Genre `json:"favourite_genres" validate:"required,dive"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "user_id", Value: userId}}
		update := bson.M{
			"$set": bson.M{
				"favourite_genres": req.FavouriteGenres,
				"updated_at":       time.Now(),
			},
		}

		result, err := database.OpenCollection("users", client).UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Favourite genres updated"})
	}
}
