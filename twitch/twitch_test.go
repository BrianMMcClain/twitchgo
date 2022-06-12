package twitch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserByLogin(t *testing.T) {
	// Reference data pulled from Twitch API docs
	expected := `{
	"data": [
			{
			"id": "141981764",
			"login": "twitchdev",
			"display_name": "TwitchDev",
			"type": "",
			"broadcaster_type": "partner",
			"description": "Supporting third-party developers building Twitch integrations from chatbots to game integrations.",
			"profile_image_url": "https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png",
			"offline_image_url": "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png",
			"view_count": 5980557,
			"email": "not-real@email.com",
			"created_at": "2016-12-14T20:32:28Z"
			}
		]
	}`

	wantID := "141981764"
	wantLogin := "twitchdev"
	wantDisplayName := "TwitchDev"
	wantViewCount := 5980557
	wantCreatedAt, _ := time.Parse("2006-01-02T15:04:05Z", "2016-12-14T20:32:28Z")

	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
	defer svr.Close()

	// Make the request with a mock config
	configJSON := "{\"client_id\": \"MyID\", \"client_secret\": \"MySecret\"}"
	c, _ := ParseConfig(configJSON)
	twitchConn := NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	u := twitchConn.GetUserByLogin("testUser")

	// Verify test
	if u.ID != wantID {
		t.Fatalf(`GetUserByLogin(login) = got %s, want %s`, u.ID, wantID)
	} else if u.Login != wantLogin {
		t.Fatalf(`GetUserByLogin(login) = got %s, want %s`, u.Login, wantLogin)
	} else if u.DisplayName != wantDisplayName {
		t.Fatalf(`GetUserByLogin(login) = got %s, want %s`, u.DisplayName, wantDisplayName)
	} else if u.ViewCount != wantViewCount {
		t.Fatalf(`GetUserByLogin(login) = got %d, want %d`, u.ViewCount, wantViewCount)
	} else if u.CreatedAt != wantCreatedAt {
		t.Fatalf(`GetUserByLogin(login) = got %s, want %s`, u.CreatedAt, wantCreatedAt)
	}

}

func TestGetLoggedInUserCached(t *testing.T) {
	// Reference data pulled from Twitch API docs
	testCreatedAtTime, _ := time.Parse("2006-01-02T15:04:05Z", "2016-12-14T20:32:28Z")
	u1 := User{
		"141981764",
		"twitchdev",
		"TwitchDev",
		"partner",
		"Supporting third-party developers building Twitch integrations from chatbots to game integrations.",
		"https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png",
		"https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png",
		5980557,
		testCreatedAtTime,
	}

	configJSON := "{\"client_id\": \"MyID\", \"client_secret\": \"MySecret\"}"
	c, _ := ParseConfig(configJSON)
	twitchConn := NewTwitch(c)
	twitchConn.user = u1
	u2 := twitchConn.GetLoggedInUser()

	if u1 != u2 {
		t.Fatalf(`GetLoggedInUser() = got %v, want cached %v`, u2, u1)
	}
}

func TestGetLoggedInUserNotCached(t *testing.T) {
	// Reference data pulled from Twitch API docs
	testCreatedAtTime, _ := time.Parse("2006-01-02T15:04:05Z", "2016-12-14T20:32:28Z")
	expected := User{
		"141981764",
		"twitchdev",
		"TwitchDev",
		"partner",
		"Supporting third-party developers building Twitch integrations from chatbots to game integrations.",
		"https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png",
		"https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png",
		5980557,
		testCreatedAtTime,
	}
	restData := `{
		"data": [
				{
				"id": "141981764",
				"login": "twitchdev",
				"display_name": "TwitchDev",
				"type": "",
				"broadcaster_type": "partner",
				"description": "Supporting third-party developers building Twitch integrations from chatbots to game integrations.",
				"profile_image_url": "https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png",
				"offline_image_url": "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png",
				"view_count": 5980557,
				"email": "not-real@email.com",
				"created_at": "2016-12-14T20:32:28Z"
				}
			]
		}`

	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, restData)
		}))
	defer svr.Close()

	// Make the request with a mock config
	configJSON := "{\"client_id\": \"MyID\", \"client_secret\": \"MySecret\"}"
	c, _ := ParseConfig(configJSON)
	twitchConn := NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	u := twitchConn.GetLoggedInUser()

	// Verify tests
	if expected != u {
		t.Fatalf(`GetLoggedInUser() = got %v, want %v`, expected, u)
	}
}
