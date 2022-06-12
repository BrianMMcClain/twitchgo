package twitch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testConfigJSON = "{\"client_id\": \"MyID\", \"client_secret\": \"MySecret\"}"
var testUserJSON = `{
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
var testCreatedAtTime, _ = time.Parse("2006-01-02T15:04:05Z", "2016-12-14T20:32:28Z")
var testUser = User{
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
var testStreamsList = `{
	"data": [{
		"id": "141981764",
		"user_id": "141981764",
		"user_login": "twitchdev",
		"user_name": "TwitchDev",
		"game_id": "509658",
		"game_name": "Just Chatting",
		"type": "live",
		"title": "Welcome To Twitch!",
		"viewer_count": 111111,
		"started_at": "2022-06-12T03:55:05Z",
		"language": "en",
		"thumbnail_url": "https://localhost/preview.jpg"
	}, {
		"id": "141981765",
		"user_id": "141981765",
		"user_login": "twitchdev2",
		"user_name": "TwitchDev2",
		"game_id": "509658",
		"game_name": "Just Chatting",
		"type": "live",
		"title": "Welcome To Twitch!",
		"viewer_count": 111111,
		"started_at": "2022-06-11T18:58:01Z",
		"language": "en",
		"thumbnail_url": "https://localhost/preview.jpg"
	}],
	"pagination": {}
}`

func TestGetUserByLogin(t *testing.T) {
	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testUserJSON)
		}))
	defer svr.Close()

	// Make the request with a mock config
	c, _ := ParseConfig(testConfigJSON)
	twitchConn := NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	u := twitchConn.GetUserByLogin("testUser")

	// Verify test
	wantID := "141981764"
	wantLogin := "twitchdev"
	wantDisplayName := "TwitchDev"
	wantViewCount := 5980557
	wantCreatedAt, _ := time.Parse("2006-01-02T15:04:05Z", "2016-12-14T20:32:28Z")

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
	c, _ := ParseConfig(testConfigJSON)
	twitchConn := NewTwitch(c)
	twitchConn.user = testUser
	u2 := twitchConn.GetLoggedInUser()

	if testUser != u2 {
		t.Fatalf(`GetLoggedInUser() = got %v, want cached %v`, u2, testUser)
	}
}

func TestGetLoggedInUserNotCached(t *testing.T) {
	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testUserJSON)
		}))
	defer svr.Close()

	// Make the request with a mock config
	c, _ := ParseConfig(testConfigJSON)
	twitchConn := NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	u := twitchConn.GetLoggedInUser()

	// Verify tests
	if testUser != u {
		t.Fatalf(`GetLoggedInUser() = got %v, want %v`, u, testUser)
	}
}

func TestGetFollowedStreams(t *testing.T) {
	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testStreamsList)
		}))
	defer svr.Close()

	// Make the request with a mock config

	c, _ := ParseConfig(testConfigJSON)
	twitchConn := NewTwitch(c)
	twitchConn.user = testUser
	twitchConn.BaseApiUrl = svr.URL
	streams := twitchConn.GetFollowedStreams()

	// Verify tests
	wantStreamCount := 2
	wantStream1ID := "141981764"
	wantStream2ID := "141981765"
	wantStream1Login := "twitchdev"
	wantStream2Login := "twitchdev2"
	wantDisplayName := "TwitchDev"
	wantGameId := "509658"
	wantGameName := "Just Chatting"
	wantType := "live"
	wantTitle := "Welcome To Twitch!"
	wantViewerCount := 111111
	wantStartedAt, _ := time.Parse("2006-01-02T15:04:05Z", "2022-06-12T03:55:05Z")
	wantThumbnailURL := "https://localhost/preview.jpg"

	if len(streams) != wantStreamCount {
		t.Fatalf(`GetFollowedStreams() = got %d, want %d`, len(streams), wantStreamCount)
	} else if streams[0].ID != wantStream1ID {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].ID, wantStream1ID)
	} else if streams[1].ID != wantStream2ID {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[1].ID, wantStream2ID)
	} else if streams[0].UserLogin != wantStream1Login {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].UserLogin, wantStream1Login)
	} else if streams[1].UserLogin != wantStream2Login {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[1].UserLogin, wantStream2Login)
	} else if streams[0].UserName != wantDisplayName {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].UserName, wantDisplayName)
	} else if streams[0].GameID != wantGameId {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].GameID, wantGameId)
	} else if streams[0].GameName != wantGameName {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].GameName, wantGameName)
	} else if streams[0].Type != wantType {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].Type, wantType)
	} else if streams[0].Title != wantTitle {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].Title, wantTitle)
	} else if streams[0].ViewerCount != wantViewerCount {
		t.Fatalf(`GetFollowedStreams() = got %d, want %d`, streams[0].ViewerCount, wantViewerCount)
	} else if streams[0].StartedAt != wantStartedAt {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].StartedAt, wantStartedAt)
	} else if streams[0].ThumbnailURL != wantThumbnailURL {
		t.Fatalf(`GetFollowedStreams() = got %s, want %s`, streams[0].ThumbnailURL, wantThumbnailURL)
	}
}
