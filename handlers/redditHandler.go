package handlers

import (
	"backend/services"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func TrendingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := services.FetchRedditTrendingPosts()
	if err != nil {
		log.Printf("Failed to fetch trending posts: %v", err)
		http.Error(w, "Failed to fetch trending posts", http.StatusInternalServerError)
		return
	}
	response := Response{
		Status:  "success",
		Message: "Trending posts fetched successfully",
		Data:    posts,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func FetchTrendingInDB(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	posts, err := services.RetrieveRedditData(collection)
	if err != nil {
		log.Printf("Failed to retrieve trending posts from DB: %v", err)
		http.Error(w, "Failed to retrieve trending posts from DB", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{Status: "success", Data: posts})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
