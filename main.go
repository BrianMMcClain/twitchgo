package main

import (
	"fmt"

	twitchbuddy "github.com/brianmmcclain/twitch-buddy-go/lib"
)

func main() {
	config := twitchbuddy.LoadConfig("config/config.json")
	twitch := twitchbuddy.NewTwitch(config)
	twitch.Auth()
	u := twitch.GetUser()
	fmt.Printf("Hello, %s!\n", u.DisplayName)
	//config.WriteConfig("config/config.json")
}
