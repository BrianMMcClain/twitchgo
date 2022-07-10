package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/brianmmcclain/twitchgo"
)

func main() {
	config := twitchgo.LoadConfig("config/config.json")
	twitchConn := twitchgo.NewTwitch(config)
	twitchConn.Auth()

	u := twitchConn.GetLoggedInUser()
	fmt.Printf("Hello, %s!\n", u.DisplayName)

	streams := twitchConn.GetFollowedStreams(u)
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

	twitchConn.ChatConnect(streams[in-1].UserLogin, chatHandler)

	fmt.Scanln()
}

func chatHandler(m *twitchgo.Message) {
	fmt.Printf("%s: %s\n", m.Sender, m.Text)
}
