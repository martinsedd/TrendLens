package handlers

import (
	"backend/services"
	"encoding/json"
	"net/http"
)

func TrendingTopicsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topics, _ := services.FetchRedditTrendingTopics()

	if err := json.NewEncoder(w).Encode(topics); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
