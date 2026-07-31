package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	defaultPage  = 1
	defaultLimit = 12
	maxLimit     = 50
)

// parsePagination reads the page/limit query params, clamping limit to a
// sane range so nobody can ask for the whole collection in one go.
func parsePagination(c *gin.Context) (page, limit int64, err error) {
	page = defaultPage
	limit = defaultLimit

	if raw := c.Query("page"); raw != "" {
		page, err = strconv.ParseInt(raw, 10, 64)
		if err != nil || page < 1 {
			return 0, 0, errors.New("page must be a positive integer")
		}
	}
	if raw := c.Query("limit"); raw != "" {
		limit, err = strconv.ParseInt(raw, 10, 64)
		if err != nil || limit < 1 || limit > maxLimit {
			return 0, 0, errors.New("limit must be between 1 and 50")
		}
	}
	return page, limit, nil
}

// buildMovieFilter turns the query params into a Mongo filter.
// q matches against the title (case-insensitive), genre matches the
// genre_name subdocument.
func buildMovieFilter(query, genre string) bson.M {
	filter := bson.M{}
	if query != "" {
		filter["title"] = bson.M{"$regex": query, "$options": "i"}
	}
	if genre != "" {
		filter["genre.genre_name"] = genre
	}
	return filter
}

// sortField whitelists the fields we allow sorting on.
func sortField(raw string) string {
	switch raw {
	case "title", "imdb_id":
		return raw
	case "rating", "ranking":
		return "ranking.ranking_value"
	default:
		return "title"
	}
}
