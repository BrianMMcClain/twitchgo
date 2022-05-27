package lib

import "time"

type UserResponse struct {
	Data []User `json:"data"`
}

type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	CreatedAt       time.Time `json:"created_at"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	ViewCount       int       `json:"view_count"`
}
