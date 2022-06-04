package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Twitch struct {
	config    *Configuration
	server    http.Server
	user_id   string
	waitGroup *sync.WaitGroup
}

func NewTwitch(config *Configuration) *Twitch {
	t := new(Twitch)
	t.config = config
	return t
}

func (t *Twitch) Auth() {

	if len(t.config.Auth) == 0 {
		t.fetchAuthCode()
	}

	if len(t.config.Token.RefreshToken) > 0 && !t.config.Token.Expires.Before(time.Now()) {
		log.Println("Token still valid, reusing")
	} else {
		log.Println("Token expired, refreshing")
		token := t.fetchToken()
		t.config.Token = *token
		t.config.WriteConfig(t.config.path)
	}
}

func (t *Twitch) fetchAuthCode() {
	// Configure the waitgroup
	t.waitGroup = &sync.WaitGroup{}
	t.waitGroup.Add(1)

	// Setup the local HTTP server
	t.server = http.Server{Addr: ":8080"}
	http.HandleFunc("/", t.authCallback)
	go func() {
		err := t.server.ListenAndServe()
		log.Println("Server listening on port 8080")
		if err != http.ErrServerClosed {
			log.Fatalf("Error starting local HTTP server: %v", err)
		}
	}()

	// Open the authentication URL to get an auth token for the logged in user
	authURL := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&redirect_uri=http://localhost:8080&client_id=%s&scope=user%%3Aread%%3Afollows", t.config.ClientID)

	log.Printf("Please authenticate using your browser: %s\n", authURL)

	// Once we receive a response, shutdown the http server
	t.waitGroup.Wait()
	t.server.Shutdown(context.TODO())
}

func (t *Twitch) authCallback(w http.ResponseWriter, req *http.Request) {
	// Store the auth token
	code := req.URL.Query().Get("code")
	t.config.Auth = code

	resp := "Logged in! You may now close the window."
	w.Write([]byte(resp))

	// Notify the waitgroup that we've receive the token and to shutdown the http server
	t.waitGroup.Done()
}

func (t *Twitch) fetchToken() *Token {
	// Build the request, providing the access code
	tokenURL := "https://id.twitch.tv/oauth2/token"
	data := url.Values{
		"client_id":     {t.config.ClientID},
		"client_secret": {t.config.ClientSecret},
		"code":          {t.config.Auth},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8080"},
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading token response: %v", err)
	}

	// Unmarshal the token
	token := new(Token)
	json.Unmarshal(respBody, &token)
	token.Expires = time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
	return token
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

func (t *Twitch) GetLoggedInUser() User {
	requestURL := "https://api.twitch.tv/helix/users"
	respBody := sendRequest(requestURL, t)
	u := new(UserResponse)
	json.Unmarshal(respBody, &u)
	t.user_id = u.Data[0].ID
	return u.Data[0]
}

func (t *Twitch) GetFollowedStreams() []Stream {

	// Get the logged in user's ID if not already cached
	if len(t.user_id) == 0 {
		t.GetLoggedInUser()
	}

	requestURL := fmt.Sprintf("https://api.twitch.tv/helix/streams/followed?user_id=%s", t.user_id)
	respBody := sendRequest(requestURL, t)
	streams := new(StreamsResponse)
	json.Unmarshal(respBody, &streams)
	return streams.Data
}
