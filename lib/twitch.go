package lib

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
)

type Twitch struct {
	config    *Configuration
	token     RefreshToken
	server    http.Server
	waitGroup *sync.WaitGroup
}

type RefreshToken struct {
	token string
}

func NewTwitch(config *Configuration) *Twitch {
	t := new(Twitch)
	t.config = config
	return t
}

func (t *Twitch) Auth() *RefreshToken {

	t.fetchAuthToken()

	r := new(RefreshToken)
	return r
}

func (t *Twitch) fetchAuthToken() {
	// Configure the waitgroup
	t.waitGroup = &sync.WaitGroup{}
	t.waitGroup.Add(1)

	// Setup the local HTTP server
	t.server = http.Server{Addr: ":8080"}
	http.HandleFunc("/", t.authCallback)
	go func() {
		err := t.server.ListenAndServe()
		log.Printf("Server listening: %v", t.server)
		if err != http.ErrServerClosed {
			log.Fatalf("Error starting local HTTP server: %v", err)
		}
	}()

	// Open the authentication URL to get an auth token for the logged in user
	authURL := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&redirect_uri=http://localhost:8080&client_id=%s&scope=user%%3Aread%%3Afollows", t.config.ClientID)
	log.Printf("Auth URL: %s", authURL)
	exec.Command("open", authURL).Start()

	// Once we receive a response, shutdown the http server
	t.waitGroup.Wait()
	t.server.Shutdown(context.TODO())
}

func (t *Twitch) authCallback(w http.ResponseWriter, req *http.Request) {
	// Store the auth token
	code := req.URL.Query().Get("code")
	t.config.AuthToken = code
	log.Printf("Token received: %s", t.config.AuthToken)

	//TODO: Send response

	// Notify the waitgroup that we've receive the token and to shutdown the http server
	t.waitGroup.Done()
}
