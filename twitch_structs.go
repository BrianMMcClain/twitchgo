package twitchgo

import "time"

type UserResponse struct {
	Data []User `json:"data"`
}

type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type Stream struct {
	ID           string    `json:"id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameName     string    `json:"game_name"`
	GameID       string    `json:"game_id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type StreamsResponse struct {
	Data []Stream `json:"data"`
}

type EmotesResponse struct {
	Data     []Emote `json:"data"`
	Template string  `json:"template"`
}
type Emote struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Images    EmoteImages `json:"images"`
	Tier      string      `json:"tier"`
	Format    []string    `json:"format"`
	Scale     []string    `json:"scale"`
	ThemeMode []string    `json:"theme_mode"`
	Type      string      `json:"emote_type"`
	SetID     string      `json:"emote_set_id"`
}

type EmoteImages struct {
	URL1x string `json:"url_1x"`
	URL2x string `json:"url_2x"`
	URL4x string `json:"url_4x"`
}

type ChatSettings struct {
	BroadcasterID        string `json:"broadcaster_id"`
	SlowMode             bool   `json:"slow_mode"`
	SlowModeWaitTime     int    `json:"slow_mode_wait_time"`
	FollowerMode         bool   `json:"follower_mode"`
	FollowerModeDuration int    `json:"follower_mode_duration"`
	SubscriberMode       bool   `json:"subscriber_mode"`
	EmoteMode            bool   `json:"emote_mode"`
	UniqueChatMode       bool   `json:"unique_chat_mode"`
}

type ChatSettingsResponse struct {
	Data []ChatSettings `json:"data"`
}
