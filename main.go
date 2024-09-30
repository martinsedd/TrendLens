package main

import (
	"backend/config"
	"backend/handlers"
	"backend/services"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
)

var redisClient *redis.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	config.InitializeMongoClient()

	topics, err := services.FetchRedditTrendingTopics()
	if err != nil {
		log.Fatalf("Error fetching Reddit trending topics: %v", err)
	}

	err = services.StoreRedditPosts(topics)
	if err != nil {
		log.Fatalf("Error storing Reddit posts: %v", err)
	}

	log.Println("Reddit trending topics successfully fetched and stored in MongoDB.")

	router := mux.NewRouter()
	router.HandleFunc("/trending", handlers.TrendingTopicsHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	response, err := services.FetchRedditTrendingTopics()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
