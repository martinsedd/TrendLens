package services

import (
	"backend/config"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	RedditAuthURL   = "https://www.reddit.com/api/v1/access_token"
	RedditAPIURL    = "https://oauth.reddit.com"
	RedditUserAgent = "TrendlensBot/0.1 by Due_Effective477"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func FetchRedditAccessToken() (string, error) {
	clientID := config.GetEnv("REDDIT_CLIENT_ID", "")
	clientSecret := config.GetEnv("REDDIT_CLIENT_SECRET", "")
	username := config.GetEnv("REDDIT_USERNAME", "")
	password := config.GetEnv("REDDIT_PASSWORD", "")

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", RedditAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %v", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("User-Agent", RedditUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute auth request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return "", fmt.Errorf("failed to parse response body: %v", err)
	}

	return token.AccessToken, nil
}

func FetchRedditTrendingTopics() ([]models.TrendingTopic, error) {
	accessToken, err := FetchRedditAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", RedditAPIURL+"/r/all/hot", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reddit API request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("User-Agent", RedditUserAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send Reddit API request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Reddit API response: %v", err)
	}

	var redditResponse map[string]interface{}
	if err := json.Unmarshal(body, &redditResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Reddit API response: %v", err)
	}

	var trendingTopics []models.TrendingTopic
	posts := redditResponse["data"].(map[string]interface{})["children"].([]interface{})
	for _, post := range posts {
		data := (post.(map[string]interface{})["data"]).(map[string]interface{})
		trendingTopics = append(trendingTopics, models.TrendingTopic{
			ID:         data["id"].(string),
			Name:       data["title"].(string),
			VolumeUp:   int(data["ups"].(float64)),
			VolumeDown: int(data["downs"].(float64)),
		})
	}
	return trendingTopics, nil
}

func StoreRedditPosts(posts []models.TrendingTopic) error {
	collection := config.MongoClient.Database("trendlens").Collection("reddit_posts")

	for _, post := range posts {
		redditPost := models.RedditPost{
			Title:      post.Name,
			Upvotes:    post.VolumeUp,
			Downvotes:  post.VolumeDown,
			Subreddit:  "r/all",
			PermaLink:  "https://www.reddit.com/r/all/comments/" + post.ID,
			InsertedAt: time.Now(),
		}
		_, err := collection.InsertOne(context.Background(), redditPost)
		if err != nil {
			return fmt.Errorf("failed to insert Reddit post into MongoDB: %v", err)
		}
	}
	fmt.Println("Reddit posts stored successfully")
	return nil
}
