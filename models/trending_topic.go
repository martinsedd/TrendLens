package models

type TrendingTopic struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	VolumeUp   int    `json:"volume_up"`
	VolumeDown int    `json:"volume_down"`
}
