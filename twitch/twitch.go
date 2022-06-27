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
	config     *Configuration
	server     http.Server
	user       User
	waitGroup  *sync.WaitGroup
	BaseApiUrl string
}

func NewTwitch(config *Configuration) *Twitch {
	t := new(Twitch)
	t.config = config
	t.BaseApiUrl = "https://api.twitch.tv/helix"
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
	requestURL := fmt.Sprintf("%s/users", t.BaseApiUrl)
	if len(login) > 0 {
		requestURL += "?login=" + login
	}
	respBody := sendRequest(requestURL, t)
	u := new(UserResponse)
	json.Unmarshal(respBody, &u)
	return u.Data[0]
}

func (t *Twitch) GetLoggedInUser() User {
	if len(t.user.ID) == 0 {
		u := t.GetUserByLogin(t.user.Login)
		t.user = u
		return u
	} else {
		return t.user
	}
}

func (t *Twitch) GetFollowedStreams(u User) []Stream {
	requestURL := fmt.Sprintf("%s/streams/followed?user_id=%s", t.BaseApiUrl, u.ID)
	respBody := sendRequest(requestURL, t)
	streams := new(StreamsResponse)
	json.Unmarshal(respBody, &streams)
	return streams.Data
}

func (t *Twitch) GetChannelEmotes(u User) ([]Emote, string) {
	requestURL := fmt.Sprintf("%s/chat/emotes?broadcaster_id=%s", t.BaseApiUrl, u.ID)
	respBody := sendRequest(requestURL, t)
	emotes := new(EmotesResponse)
	json.Unmarshal(respBody, &emotes)
	return emotes.Data, emotes.Template
}
