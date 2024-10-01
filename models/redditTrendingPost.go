package models

// TrendingPost represents a structure for a trending post in the application.
// It includes fields that capture the essential information about the trending post,
// such as its ID, name, and volume of votes.
type TrendingPost struct {
	ID         string `json:"id"`          // Unique identifier for the trending post
	Name       string `json:"name"`        // Name or title of the trending post
	VolumeUp   int    `json:"volume_up"`   // Number of upvotes for the trending post
	VolumeDown int    `json:"volume_down"` // Number of downvotes for the trending post
}
