package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestParsePaginationDefaults(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/movies", nil)

	page, limit, err := parsePagination(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if page != 1 || limit != 12 {
		t.Errorf("expected page=1 limit=12, got page=%d limit=%d", page, limit)
	}
}

func TestParsePaginationFromQuery(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/movies?page=3&limit=20", nil)

	page, limit, err := parsePagination(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if page != 3 || limit != 20 {
		t.Errorf("expected page=3 limit=20, got page=%d limit=%d", page, limit)
	}
}

func TestParsePaginationRejectsBadInput(t *testing.T) {
	bad := []string{"page=0", "page=-1", "page=abc", "limit=0", "limit=999", "page=2&limit=-5"}

	for _, query := range bad {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/movies?"+query, nil)

		if _, _, err := parsePagination(ctx); err == nil {
			t.Errorf("expected error for query %q", query)
		}
	}
}

func TestBuildMovieFilter(t *testing.T) {
	filter := buildMovieFilter("dune", "Sci-Fi")

	if filter["genre.genre_name"] != "Sci-Fi" {
		t.Error("genre filter missing from result")
	}

	titleRegex, ok := filter["title"].(bson.M)
	if !ok {
		t.Fatal("title filter should be a regex document")
	}
	if titleRegex["$regex"] != "dune" {
		t.Errorf("expected regex for 'dune', got %v", titleRegex["$regex"])
	}
	if titleRegex["$options"] != "i" {
		t.Errorf("expected case-insensitive match, got %v", titleRegex["$options"])
	}
}

func TestBuildMovieFilterEmpty(t *testing.T) {
	if len(buildMovieFilter("", "")) != 0 {
		t.Error("empty query should produce an empty filter")
	}
}

func TestSortFieldWhitelist(t *testing.T) {
	cases := map[string]string{
		"title":   "title",
		"imdb_id": "imdb_id",
		"rating":  "ranking.ranking_value",
		"ranking": "ranking.ranking_value",
		"random":  "title",
		"":        "title",
	}

	for in, want := range cases {
		if got := sortField(in); got != want {
			t.Errorf("sortField(%q) = %q, want %q", in, got, want)
		}
	}
}
