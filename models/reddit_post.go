package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RedditPost struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Upvotes    int                `bson:"upvotes"`
	Downvotes  int                `bson:"downvotes"`
	Subreddit  string             `bson:"subreddit"`
	PermaLink  string             `bson:"perma_link"`
	URL        string             `bson:"url"`
	InsertedAt time.Time          `bson:"inserted_at"`
}
