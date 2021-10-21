package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// DB Instance func
func DBinstance() *mongo.Client {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//MongoDb := os.Getenv("MONGODB_URL")
	MongoDb := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal( err )
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal( err )
	}
	fmt.Println("Connected to MongoDB ! ")

	return client
}

// Client Database Instance
var Client *mongo.Client = DBinstance()

// OpenCollection is a function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string ) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}