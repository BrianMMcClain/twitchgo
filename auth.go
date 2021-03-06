package twitchgo

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

func (t *Twitch) Auth() {
	if len(t.config.Token.RefreshToken) > 0 && !t.config.Token.Expires.Before(time.Now()) {
		// Token is still valid
	} else {
		// Token expired
		t.fetchAuthCode()
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
		if err != http.ErrServerClosed {
			log.Fatalf("Error starting local HTTP server: %v", err)
		}
	}()

	// Open the authentication URL to get an auth token for the logged in user
	authURL := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&redirect_uri=http://localhost:8080&client_id=%s&scope=user%%3Aread%%3Afollows+chat%%3Aread", t.config.ClientID)

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
