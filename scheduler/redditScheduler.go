package scheduler

import (
	"backend/services"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// StartRedditScheduler initializes and starts a scheduler to fetch trending posts from Reddit
// and store them in the specified MongoDB collection at regular intervals.
//
// Parameters:
//   - collection: The MongoDB collection where the fetched posts will be stored.
func StartRedditScheduler(collection *mongo.Collection) {
	// Create a new scheduler that operates in UTC time zone
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule a job to run every 5 minutes
	_, err := scheduler.Every(5).Minutes().Do(func() {
		// Fetch trending posts from Reddit
		posts, err := services.FetchRedditTrendingPosts()
		if err != nil {
			log.Printf("Error fetching Reddit trending topics: %v", err)
			return
		}

		// Store the fetched posts in the specified MongoDB collection
		err = services.StoreRedditPosts(collection, posts)
		if err != nil {
			log.Printf("Error storing Reddit posts: %v", err)
			return
		}
	})

	// Check if there was an error scheduling the job
	if err != nil {
		log.Fatalf("Error scheduling Reddit job: %v", err)
	}

	// Start the scheduler asynchronously
	scheduler.StartAsync()
}
