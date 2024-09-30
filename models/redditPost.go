package models

import (
	"time"
)

type VoteHistoryEntry struct {
	Value     int       `bson:"value"`
	Timestamp time.Time `bson:"timestamp"`
}
type RedditPost struct {
	ID              string             `bson:"_id,omitempty"`
	Title           string             `bson:"title"`
	Upvotes         int                `bson:"upvotes"`
	Downvotes       int                `bson:"downvotes"`
	Subreddit       string             `bson:"subreddit"`
	PermaLink       string             `bson:"perma_link"`
	URL             string             `bson:"url"`
	InsertedAt      time.Time          `bson:"inserted_at"`
	UpvoteHistory   []VoteHistoryEntry `bson:"upvote_history"`
	DownvoteHistory []VoteHistoryEntry `bson:"downvote_history"`
}
