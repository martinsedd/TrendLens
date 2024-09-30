package scheduler

import (
	"backend/services"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func StartRedditScheduler(collection *mongo.Collection) {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(5).Minutes().Do(func() {
		posts, err := services.FetchRedditTrendingPosts()
		if err != nil {
			log.Printf("Error fetching Reddit trending topics: %v", err)
			return
		}

		err = services.StoreRedditPosts(collection, posts)
		if err != nil {
			log.Printf("Error storing Reddit posts: %v", err)
			return
		}
	})

	if err != nil {
		log.Fatalf("Error scheduling Reddit job: %v", err)
	}

	scheduler.StartAsync()
}
