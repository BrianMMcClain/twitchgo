twitchgo
===

![Build Status](https://github.com/brianmmcclain/twitchgo/actions/workflows/test-and-build.yml/badge.svg)

A package to interact with the Twitch REST API and live chat

Example client in the `cmd` directory. Copy the `cmd/config/config.json.ph` file to `cmd/config/config.json` and replace the values with your own [generated client ID and secret](https://dev.twitch.tv/docs/authentication/register-app)

## Usage

```go
package main

import (
	"fmt"

	"github.com/brianmmcclain/twitchgo"
)

func main() {
    // Read the config and create a new client with it
    twitchConfig := twitchgo.LoadConfig("/path/to/config.json")
    twitchClient := twitchgo.NewTwitch(twitchConfig)
	
    // Authenticate to the Twitch API. Twitch uses your long-living
    // client ID and secret to generate short-lived tokens used to
    // interact with the Twitch APIs and authenticate in chat. If the token
    // has expired, this will provide the user a URL to use to get a new token
    twitchClient.Auth()

    // Connect a channels chat
    // This method takes a function as an argument which gets invoked
    // with each message received
    twitchClient.ChatConnect("CHANNEL_NAME", chatCallback)

    // Only for the sake of the demo so the application
    // doesn't exit immediately
    fmt.Scanln()
}

// The callback that gets called with each message received
func chatCallback(m *twitchgo.Message) {
    fmt.Printf("%s: %s\n", m.Sender, m.Text)
}
```