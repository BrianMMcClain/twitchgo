package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	twitchbuddy "github.com/brianmmcclain/twitch-buddy-go/twitchbuddy"
)

func main() {
	config := twitchbuddy.LoadConfig("config/config.json")
	twitch := twitchbuddy.NewTwitch(config)
	twitch.Auth()
	u := twitch.GetLoggedInUser()
	fmt.Printf("Hello, %s!\n", u.DisplayName)

	streams := twitch.GetFollowedStreams()
	for i, s := range streams {
		fmt.Printf("%d: %s (%d - %s): %s\n", i+1, s.UserName, s.ViewerCount, s.GameName, s.Title)
	}
	fmt.Print("Choose a chat to connect to: ")

	reader := bufio.NewReader(os.Stdin)
	inS, _, err := reader.ReadLine()
	if err != nil {
		log.Fatalf("Error reading input %s", err)
	}

	in, err := strconv.Atoi(string(inS))
	if err != nil {
		log.Fatalf("Error parsing input: %s", err)
	} else if in <= 0 || in > len(streams) {
		log.Fatalf("Invalid input: %d", in)
	}

	fmt.Printf("Connecting to %s . . .\n", streams[in-1].UserName)
	msgChannel := make(chan twitchbuddy.Message)
	twitch.ChatConnect(streams[in-1].UserLogin, msgChannel)
	for {
		msg := <-msgChannel
		fmt.Printf("%s: %s\n", msg.Sender, msg.Text)
	}
}
