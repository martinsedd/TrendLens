package main

import (
	"backend/config"
	"backend/handlers"
	"backend/scheduler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := config.InitializeMongoClient()
	collection := client.Database("trendlens").Collection("reddit_posts")

	scheduler.StartRedditScheduler(collection)

	router := mux.NewRouter()
	router.HandleFunc("/trending", handlers.TrendingHandler).Methods("GET")
	router.HandleFunc("/stored_posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.FetchTrendingInDB(w, r, collection)
	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
