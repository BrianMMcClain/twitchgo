package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Twitch struct {
	config    *Configuration
	server    http.Server
	user      User
	waitGroup *sync.WaitGroup
}

func NewTwitch(config *Configuration) *Twitch {
	t := new(Twitch)
	t.config = config
	return t
}

func sendRequest(requestURL string, t *Twitch) []byte {
	// Build the request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.config.Token.AccessToken))
	req.Header.Add("Client-Id", t.config.ClientID)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing request: %v\n", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	return respBody
}

func (t *Twitch) GetUserByLogin(login string) User {
	requestURL := "https://api.twitch.tv/helix/users"
	if len(login) > 0 {
		requestURL += "?login=" + login
	}
	respBody := sendRequest(requestURL, t)
	u := new(UserResponse)
	json.Unmarshal(respBody, &u)
	t.user = u.Data[0]
	return u.Data[0]
}

func (t *Twitch) GetLoggedInUser() User {
	if len(t.user.ID) == 0 {
		return t.GetUserByLogin(t.user.Login)
	} else {
		return t.user
	}
}

func (t *Twitch) GetFollowedStreams() []Stream {

	// Get the logged in user's ID if not already cached
	if len(t.user.ID) == 0 {
		t.GetLoggedInUser()
	}

	requestURL := fmt.Sprintf("https://api.twitch.tv/helix/streams/followed?user_id=%s", t.user.ID)
	respBody := sendRequest(requestURL, t)
	streams := new(StreamsResponse)
	json.Unmarshal(respBody, &streams)
	return streams.Data
}
