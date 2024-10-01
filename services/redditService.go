package services

import (
	"backend/config"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jonreiter/govader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	RedditAuthURL   = "https://www.reddit.com/api/v1/access_token" // URL for Reddit authentication
	RedditAPIURL    = "https://oauth.reddit.com"                   // Base URL for Reddit API
	RedditUserAgent = "TrendlensBot/0.1 by Due_Effective477"       // User agent string for Reddit requests
)

// Token represents the structure of the access token received from Reddit.
type Token struct {
	AccessToken string `json:"access_token"` // The access token for authenticating API requests
	TokenType   string `json:"token_type"`   // The type of token (typically "bearer")
	ExpiresIn   int    `json:"expires_in"`   // The expiration time of the token in seconds
}

// FetchRedditAccessToken retrieves an access token from Reddit using the password grant type.
// It returns the access token as a string or an error if the request fails.
func FetchRedditAccessToken() (string, error) {
	clientID := config.GetEnv("REDDIT_CLIENT_ID", "")
	clientSecret := config.GetEnv("REDDIT_CLIENT_SECRET", "")
	username := config.GetEnv("REDDIT_USERNAME", "")
	password := config.GetEnv("REDDIT_PASSWORD", "")

	// Prepare data for the POST request
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", RedditAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %v", err)
	}

	req.SetBasicAuth(clientID, clientSecret) // Set basic auth credentials
	req.Header.Add("User-Agent", RedditUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute auth request: %v", err)
	}
	defer res.Body.Close() // Ensure response body is closed

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return "", fmt.Errorf("failed to parse response body: %v", err)
	}

	return token.AccessToken, nil // Return the access token
}

// FetchRedditTrendingPosts retrieves the trending posts from Reddit's "hot" section.
// It returns a slice of TrendingPost models or an error if the request fails.
func FetchRedditTrendingPosts() ([]models.TrendingPost, error) {
	accessToken, err := FetchRedditAccessToken() // Get access token
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", RedditAPIURL+"/r/all/hot", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reddit API request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken) // Set the authorization header
	req.Header.Add("User-Agent", RedditUserAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send Reddit API request: %v", err)
	}
	defer res.Body.Close() // Ensure response body is closed

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Reddit API response: %v", err)
	}

	var redditResponse map[string]interface{}
	if err := json.Unmarshal(body, &redditResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Reddit API response: %v", err)
	}

	data, ok := redditResponse["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected Reddit API response format")
	}

	posts, ok := data["children"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected Reddit API response format")
	}

	var trendingPosts []models.TrendingPost
	for _, post := range posts {
		postData, ok := post.(map[string]interface{})["data"].(map[string]interface{})
		if !ok {
			continue // Skip malformed post data
		}
		// Append the trending post to the slice
		trendingPosts = append(trendingPosts, models.TrendingPost{
			ID:         postData["id"].(string),
			Name:       postData["title"].(string),
			VolumeUp:   int(postData["ups"].(float64)),
			VolumeDown: int(postData["downs"].(float64)),
		})
	}
	return trendingPosts, nil // Return the slice of trending posts
}

// StoreRedditPosts stores or updates the trending posts in the MongoDB collection.
// It performs sentiment analysis on the post titles and keeps track of voting history.
func StoreRedditPosts(collection *mongo.Collection, posts []models.TrendingPost) error {
	analyzer := govader.NewSentimentIntensityAnalyzer() // Initialize the sentiment analyzer

	for _, post := range posts {
		sentiment := analyzer.PolarityScores(post.Name) // Analyze sentiment of the post title

		// Determine the sentiment label based on the compound score
		sentimentLabel := "neutral"
		if sentiment.Compound >= 0.05 {
			sentimentLabel = "positive"
		} else if sentiment.Compound <= -0.05 {
			sentimentLabel = "negative"
		}

		filter := bson.M{"id": post.ID} // Create a filter for MongoDB query
		var existingPost models.RedditPost

		err := collection.FindOne(context.Background(), filter).Decode(&existingPost)
		if err != nil && err.Error() != "mongo: no documents in result" {
			return fmt.Errorf("error fetching Reddit post from MongoDB: %v", err)
		}

		// Prepare the update for the MongoDB document
		update := bson.M{
			"$set": bson.M{
				"title":       post.Name,
				"upvotes":     post.VolumeUp,
				"downvotes":   post.VolumeDown,
				"subreddit":   "r/all",
				"perma_link":  "https://reddit.com/r/all/comments/" + post.ID,
				"url":         "https://reddit.com/r/all/comments/" + post.ID,
				"inserted_at": time.Now(),
				"sentiment":   sentimentLabel,
			},
		}

		if existingPost.ID != "" {
			// Track upvote history if the upvotes have changed
			if existingPost.Upvotes != post.VolumeUp {
				upvoteEntry := models.VoteHistoryEntry{
					Value:     post.VolumeUp,
					Timestamp: time.Now(),
				}
				update["$push"] = bson.M{"upvote_history": upvoteEntry}
			}
			// Track downvote history if the downvotes have changed
			if existingPost.Downvotes != post.VolumeDown {
				downvoteEntry := models.VoteHistoryEntry{
					Value:     post.VolumeDown,
					Timestamp: time.Now(),
				}
				// If $push was already set, we need to merge the push for downvote_history
				if pushData, ok := update["$push"]; ok {
					pushData := pushData.(bson.M)
					pushData["downvote_history"] = downvoteEntry
					update["$push"] = pushData
				} else {
					update["$push"] = bson.M{"downvote_history": downvoteEntry}
				}
			}
		}

		// Use upsert to insert or update the document
		upsert := true
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err = collection.UpdateOne(context.Background(), filter, update, &opts)
		if err != nil {
			return fmt.Errorf("failed to upsert Reddit post into MongoDB: %v", err)
		}
	}

	fmt.Println("Reddit posts stored or updated successfully: " + time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

// RetrieveRedditData retrieves all Reddit posts from the specified MongoDB collection.
// It returns a slice of RedditPost models or an error if the retrieval fails.
func RetrieveRedditData(collection *mongo.Collection) ([]models.RedditPost, error) {
	filter := bson.M{} // Empty filter to retrieve all documents

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Reddit posts from MongoDB: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx) // Ensure the cursor is closed after usage
		if err != nil {
			fmt.Println("Failed to close cursor: ", err)
		}
	}(cursor, context.Background())

	var posts []models.RedditPost
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, fmt.Errorf("failed to decode Reddit posts from cursor: %v", err)
	}

	return posts, nil // Return the slice of retrieved posts
}
