package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return client
}

func OpenCollection(name string, client *mongo.Client) *mongo.Collection {
	if client == nil {
		log.Fatal("database client is not initialized; call Connect first")
	}

	db := os.Getenv("DATABASE_NAME")
	if db == "" {
		log.Fatal("DATABASE_NAME is not set")
	}

	return client.Database(db).Collection(name)
}
