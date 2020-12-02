package persister

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo() (*mongo.Database, func()) {

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	// connect to MongoDB
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicln(err)
	}
	disconnectMongo := func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panicln(err)
		}
	}
	db := client.Database(dbName)

	return db, disconnectMongo
}
