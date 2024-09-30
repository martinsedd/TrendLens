package handlers

import (
	"backend/models"
	"backend/services"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
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

func FetchFilteredPostsHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	sentiment := r.URL.Query().Get("sentiment")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	filter := bson.M{}
	if sentiment != "" {
		filter["sentiment"] = sentiment
	}

	skip := int64((page - 1) * limit)

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(skip)

	cursor, err := collection.Find(r.Context(), filter, findOptions)
	if err != nil {
		log.Printf("Failed to find posts: %v", err)
		http.Error(w, "Failed to find posts", http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Failed to close cursor: %v", err)
		}
	}(cursor, r.Context())

	var posts []models.RedditPost
	if err := cursor.All(context.Background(), &posts); err != nil {
		log.Printf("Failed to decode posts: %v", err)
		http.Error(w, "Failed to decode posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Error encoding posts to JSON: %v", err)
		http.Error(w, "Error encoding posts to JSON", http.StatusInternalServerError)
		return
	}
}
