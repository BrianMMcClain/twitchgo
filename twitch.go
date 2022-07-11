package twitchgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func sendRequest(requestURL string, t *Twitch) ([]byte, error) {
	// Build the request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.config.Token.AccessToken))
	req.Header.Add("Client-Id", t.config.ClientID)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error performing request: %v", err))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading response body: %v", err))
	}
	return respBody, nil
}

func (t *Twitch) GetUserByLogin(login string) (User, error) {
	requestURL := fmt.Sprintf("%s/users", t.BaseApiUrl)
	if len(login) > 0 {
		requestURL += "?login=" + login
	}
	respBody, err := sendRequest(requestURL, t)
	if err != nil {
		return User{}, err
	}
	u := new(UserResponse)
	json.Unmarshal(respBody, &u)
	return u.Data[0], nil
}

func (t *Twitch) GetLoggedInUser() (User, error) {
	if len(t.user.ID) == 0 {
		u, err := t.GetUserByLogin(t.user.Login)
		if err != nil {
			return User{}, err
		}
		t.user = u
		return u, nil
	} else {
		return t.user, nil
	}
}

func (t *Twitch) GetFollowedStreams(u User) ([]Stream, error) {
	requestURL := fmt.Sprintf("%s/streams/followed?user_id=%s", t.BaseApiUrl, u.ID)
	respBody, err := sendRequest(requestURL, t)
	if err != nil {
		return nil, err
	}
	streams := new(StreamsResponse)
	json.Unmarshal(respBody, &streams)
	return streams.Data, nil
}

func (t *Twitch) GetChannelEmotes(u User) ([]Emote, error) {
	requestURL := fmt.Sprintf("%s/chat/emotes?broadcaster_id=%s", t.BaseApiUrl, u.ID)
	respBody, err := sendRequest(requestURL, t)
	if err != nil {
		return nil, err
	}
	emotes := new(EmotesResponse)
	json.Unmarshal(respBody, &emotes)
	return emotes.Data, nil
}

func (t *Twitch) GetChatSettings(u User) (*ChatSettings, error) {
	requestURL := fmt.Sprintf("%s/chat/settings?broadcaster_id=%s", t.BaseApiUrl, u.ID)
	respBody, err := sendRequest(requestURL, t)
	if err != nil {
		return nil, err
	}
	settings := new(ChatSettingsResponse)
	json.Unmarshal(respBody, &settings)

	if len(settings.Data) > 0 {
		return &settings.Data[0], nil
	} else {
		return nil, errors.New("no chat settings could be retrieved")
	}
}
