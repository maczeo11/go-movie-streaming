package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/maczeo11/go-movie-streaming/Server/MagicMovieStream/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// seed replays the JSON files in ./seed into MongoDB. It's idempotent:
// each collection is wiped and re-inserted, so running it twice is safe.
// Usage: go run ./cmd/seed [path-to-seed-dir]
func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	client := database.Connect()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		log.Fatal("DATABASE_NAME is not set")
	}
	db := client.Database(dbName)

	seedDir := "seed"
	if len(os.Args) > 1 {
		seedDir = os.Args[1]
	}

	seedCollection(ctx, db, seedDir, "genres.json", "genres")
	seedCollection(ctx, db, seedDir, "rankings.json", "rankings")
	seedCollection(ctx, db, seedDir, "movies.json", "movies")

	slog.Info("seed complete")
}

func seedCollection(ctx context.Context, db *mongo.Database, dir, file, collection string) {
	var items []bson.M

	data, err := os.ReadFile(filepath.Join(dir, file))
	if err != nil {
		log.Fatalf("failed to read %s: %v", file, err)
	}
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatalf("failed to parse %s: %v", file, err)
	}

	col := db.Collection(collection)
	if _, err := col.DeleteMany(ctx, bson.D{}); err != nil {
		log.Fatalf("failed to wipe %s: %v", collection, err)
	}

	docs := make([]any, len(items))
	for i := range items {
		docs[i] = items[i]
	}

	res, err := col.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("failed to seed %s: %v", collection, err)
	}
	slog.Info("seeded collection", "collection", collection, "count", len(res.InsertedIDs))
}
