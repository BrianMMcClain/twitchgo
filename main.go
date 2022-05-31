package main

import (
	"fmt"

	twitchbuddy "github.com/brianmmcclain/twitch-buddy-go/lib"
)

func main() {
	config := twitchbuddy.LoadConfig("config/config.json")
	twitch := twitchbuddy.NewTwitch(config)
	twitch.Auth()
	u := twitch.GetLoggedInUser()
	fmt.Printf("Hello, %s!\n", u.DisplayName)

	streams := twitch.GetFollowedStreams()
	for _, s := range streams {
		fmt.Printf("%s (%d - %s): %s\n", s.UserName, s.ViewerCount, s.GameName, s.Title)
	}
}
