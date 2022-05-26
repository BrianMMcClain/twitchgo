package main

import (
	"log"

	twitchbuddy "github.com/brianmmcclain/twitch-buddy-go/lib"
)

func main() {
	config := twitchbuddy.NewConfig("config/config.json")
	log.Printf("Client ID: %s", config.ClientID)
}
