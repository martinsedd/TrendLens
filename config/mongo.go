package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// MongoClient is a global variable that holds the MongoDB client instance.
var MongoClient *mongo.Client

// InitializeMongoClient initializes a new MongoDB client connection.
//
// It retrieves the MongoDB URI from the environment variable MONGO_URI, or defaults to "mongodb://localhost:27017".
// The function sets the connection options, connects to the database, and pings the server to check connectivity.
// If successful, it assigns the client to the global MongoClient variable and returns the client instance.
//
// Returns:
//   - A pointer to the mongo.Client instance.
func InitializeMongoClient() *mongo.Client {
	// Retrieve the MongoDB URI from environment variables or use a default value
	uri := GetEnv("MONGO_URI", "mongodb://localhost:27017")
	// Set up client options with a maximum pool size, minimum pool size, and maximum connection idle time
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(20).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(10 * time.Minute)

	// Attempt to connect to the MongoDB server
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Set a timeout context for the ping operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	// Assign the connected client to the global MongoClient variable
	MongoClient = client
	return client
}

// DisconnectMongoClient disconnects the MongoDB client from the server if it is connected.
//
// It logs an error if the disconnection fails and prints a message upon successful disconnection.
func DisconnectMongoClient() {
	// Check if the MongoClient is initialized
	if MongoClient != nil {
		// Attempt to disconnect the MongoDB client
		err := MongoClient.Disconnect(context.TODO())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
		fmt.Println("Disconnected from MongoDB!")
	}
}
