package twitchgo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianmmcclain/twitchgo"
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
var testUser = twitchgo.User{
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

var testChatSettingsJson = `{
	"data": [
	  {
		"broadcaster_id": "713936733",
		"slow_mode": false,
		"slow_mode_wait_time": null,
		"follower_mode": true,
		"follower_mode_duration": 20,
		"subscriber_mode": false,
		"emote_mode": false,
		"unique_chat_mode": false,
		"non_moderator_chat_delay": true,
		"non_moderator_chat_delay_duration": 4
	  }
	]
  }`

var testChannelEmotesJson = `{
	"data": [{
		"id": "000001",
		"name": "emote1",
		"images": {
			"url_1x": "https://static-cdn.jtvnw.net/emoticons/v2/000001/static/light/1.0",
			"url_2x": "https://static-cdn.jtvnw.net/emoticons/v2/000001/static/light/2.0",
			"url_4x": "https://static-cdn.jtvnw.net/emoticons/v2/000001/static/light/3.0"
		},
		"tier": "2000",
		"emote_type": "subscriptions",
		"emote_set_id": "000001",
		"format": ["static"],
		"scale": ["1.0", "2.0", "3.0"],
		"theme_mode": ["light", "dark"]
	}, {
		"id": "000002",
		"name": "emote2",
		"images": {
			"url_1x": "https://static-cdn.jtvnw.net/emoticons/v2/000002/static/light/1.0",
			"url_2x": "https://static-cdn.jtvnw.net/emoticons/v2/000002/static/light/2.0",
			"url_4x": "https://static-cdn.jtvnw.net/emoticons/v2/000002/static/light/3.0"
		},
		"tier": "1000",
		"emote_type": "subscriptions",
		"emote_set_id": "000001",
		"format": ["static"],
		"scale": ["1.0", "2.0", "3.0"],
		"theme_mode": ["light", "dark"]
	}],
	"template": "https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}"
}`

func TestGetUserByLogin(t *testing.T) {
	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testUserJSON)
		}))
	defer svr.Close()

	// Make the request with a mock config
	c, _ := twitchgo.ParseConfig(testConfigJSON)
	twitchConn := twitchgo.NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	u, _ := twitchConn.GetUserByLogin("testUser")

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

func TestGetLoggedInUser(t *testing.T) {
	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testUserJSON)
		}))
	defer svr.Close()

	// Make the request with a mock config
	c, _ := twitchgo.ParseConfig(testConfigJSON)
	twitchConn := twitchgo.NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	u, _ := twitchConn.GetLoggedInUser()

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
	c, _ := twitchgo.ParseConfig(testConfigJSON)
	twitchConn := twitchgo.NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	streams, _ := twitchConn.GetFollowedStreams(testUser)

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

func TestGetChatSettings(t *testing.T) {
	const TEST_NAME = "GetChatSettings"

	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testChatSettingsJson)
		}))
	defer svr.Close()

	// Make the request with a mock config
	c, _ := twitchgo.ParseConfig(testConfigJSON)
	twitchConn := twitchgo.NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	settings, err := twitchConn.GetChatSettings(testUser)

	// Ensure we didn't get an error parsing the response
	verify(err, nil, TEST_NAME, "ParseSettings", t)

	// Verify tests
	wantBroadcasterID := "713936733"
	wantSlowMode := false
	wantSlowModeWaitTime := 0
	wantFollowerMode := true
	wantFollowerModeDuration := 20

	verify(wantBroadcasterID, settings.BroadcasterID, TEST_NAME, "BroadcasterID", t)
	verify(wantSlowMode, settings.SlowMode, TEST_NAME, "SlowMode", t)
	verify(wantSlowModeWaitTime, settings.SlowModeWaitTime, TEST_NAME, "SlowModeDuration", t)
	verify(wantFollowerMode, settings.FollowerMode, TEST_NAME, "FollowerMode", t)
	verify(wantFollowerModeDuration, settings.FollowerModeDuration, TEST_NAME, "FollowerModeDuration", t)
}

func TestGetChannelEmotes(t *testing.T) {
	const TEST_NAME = "GetChannelEmotes"

	// Set up the test server
	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testChannelEmotesJson)
		}))
	defer svr.Close()

	// Make the request with a mock config
	c, _ := twitchgo.ParseConfig(testConfigJSON)
	twitchConn := twitchgo.NewTwitch(c)
	twitchConn.BaseApiUrl = svr.URL
	emotes, err := twitchConn.GetChannelEmotes(testUser)

	// Verify that we didn't get an error
	verify(err, nil, TEST_NAME, "ParseEmotes", t)

	// Verify tests
	wantEmoteCount := 2
	wantEmote1Name := "emote1"
	wantEmote2Name := "emote2"
	wantEmote1Tier := "2000"
	wantEmote2Tier := "1000"

	verify(len(emotes), wantEmoteCount, TEST_NAME, "EmoteCount", t)
	verify(emotes[0].Name, wantEmote1Name, TEST_NAME, "Emote1Name", t)
	verify(emotes[1].Name, wantEmote2Name, TEST_NAME, "Emote2Name", t)
	verify(emotes[0].Tier, wantEmote1Tier, TEST_NAME, "Emote1Tier", t)
	verify(emotes[1].Tier, wantEmote2Tier, TEST_NAME, "Emote2Tier", t)
}

func verify(want, got interface{}, testName string, caseName string, t *testing.T) {
	if want != got {
		t.Fatalf(`%s() %s = got %s, want %s`, testName, caseName, want, got)
	}
}
