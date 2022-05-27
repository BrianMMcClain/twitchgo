package main

import (
	twitchbuddy "github.com/brianmmcclain/twitch-buddy-go/lib"
)

func main() {
	config := twitchbuddy.LoadConfig("config/config.json")
	twitch := twitchbuddy.NewTwitch(config)
	twitch.Auth()
	config.WriteConfig("config/config.json")
}
