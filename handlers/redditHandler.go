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

// Response represents the structure of the JSON response returned by the API.
type Response struct {
	Status  string      `json:"status"`         // Status of the response (e.g., success, error)
	Message string      `json:"message"`        // Message providing more details
	Data    interface{} `json:"data,omitempty"` // Optional data payload
}

// TrendingHandler handles the request for fetching trending posts from Reddit.
func TrendingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch trending posts from the Reddit service
	posts, err := services.FetchRedditTrendingPosts()
	if err != nil {
		log.Printf("Failed to fetch trending posts: %v", err)
		http.Error(w, "Failed to fetch trending posts", http.StatusInternalServerError)
		return
	}

	// Construct a successful response
	response := Response{
		Status:  "success",
		Message: "Trending posts fetched successfully",
		Data:    posts,
	}

	// Encode the response to JSON and send it back to the client
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// FetchTrendingInDB retrieves trending posts from the MongoDB collection.
func FetchTrendingInDB(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Retrieve trending posts from the database
	posts, err := services.RetrieveRedditData(collection)
	if err != nil {
		log.Printf("Failed to retrieve trending posts from DB: %v", err)
		http.Error(w, "Failed to retrieve trending posts from DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Encode the response with the retrieved posts
	err = json.NewEncoder(w).Encode(Response{Status: "success", Data: posts})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// FetchFilteredPostsHandler handles requests to fetch posts based on filters like sentiment, pagination, etc.
func FetchFilteredPostsHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Retrieve the sentiment filter from the query parameters
	sentiment := r.URL.Query().Get("sentiment")

	// Parse the limit parameter for pagination
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit if parsing fails or limit is invalid
	}

	// Parse the page parameter for pagination
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Default page if parsing fails or page is invalid
	}

	// Construct the filter based on sentiment
	filter := bson.M{}
	if sentiment != "" {
		filter["sentiment"] = sentiment
	}

	// Calculate the number of documents to skip for pagination
	skip := int64((page - 1) * limit)

	// Set up find options for limit and skip
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(skip)

	// Query the MongoDB collection for posts
	cursor, err := collection.Find(r.Context(), filter, findOptions)
	if err != nil {
		log.Printf("Failed to find posts: %v", err)
		http.Error(w, "Failed to find posts", http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		// Close the cursor after the function completes
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Failed to close cursor: %v", err)
		}
	}(cursor, r.Context())

	// Decode the retrieved posts into a slice
	var posts []models.RedditPost
	if err := cursor.All(context.Background(), &posts); err != nil {
		log.Printf("Failed to decode posts: %v", err)
		http.Error(w, "Failed to decode posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Encode and send the retrieved posts as a JSON response
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Error encoding posts to JSON: %v", err)
		http.Error(w, "Error encoding posts to JSON", http.StatusInternalServerError)
		return
	}
}
