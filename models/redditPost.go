package models

import (
	"time"
)

// VoteHistoryEntry represents a record of a single vote on a Reddit post.
// It contains the value of the vote (upvote or downvote) and the timestamp of when the vote was cast.
type VoteHistoryEntry struct {
	Value     int       `bson:"value"`     // The value of the vote (e.g., 1 for upvote, -1 for downvote)
	Timestamp time.Time `bson:"timestamp"` // The time when the vote was recorded
}

// RedditPost represents the structure of a Reddit post in the database.
// It includes various fields relevant to a Reddit post, such as its title, vote counts, and history of votes.
type RedditPost struct {
	ID              string             `bson:"_id,omitempty"`    // Unique identifier for the post (auto-generated if omitted)
	Title           string             `bson:"title"`            // The title of the Reddit post
	Upvotes         int                `bson:"upvotes"`          // Total number of upvotes for the post
	Downvotes       int                `bson:"downvotes"`        // Total number of downvotes for the post
	Subreddit       string             `bson:"subreddit"`        // The subreddit where the post was made
	PermaLink       string             `bson:"perma_link"`       // Permanent link to the post on Reddit
	URL             string             `bson:"url"`              // URL of the post or associated content
	InsertedAt      time.Time          `bson:"inserted_at"`      // Timestamp of when the post was inserted into the database
	UpvoteHistory   []VoteHistoryEntry `bson:"upvote_history"`   // History of upvotes on the post
	DownvoteHistory []VoteHistoryEntry `bson:"downvote_history"` // History of downvotes on the post
}
