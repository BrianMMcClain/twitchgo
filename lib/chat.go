package lib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Chat struct {
	Conn      net.Conn
	Channel   string
	Connected bool
	Joined    bool
	Twitch    *Twitch
}

type Message struct {
	Sender string
	Text   string
}

func (t *Twitch) ChatConnect(channel string) {
	CHAT_HOST := "irc.chat.twitch.tv:6667"

	// Build the chat struct
	chat := new(Chat)
	chat.Channel = channel
	chat.Twitch = t

	// Connect to the server
	conn, err := net.Dial("tcp", CHAT_HOST)
	if err != nil {
		log.Fatalf("Could not connect to chat server: %s", err)
	}
	chat.Conn = conn

	// Authenticate
	// TODO: Keep things simple for now, only worry about basic chat
	// chat.sendMsg("CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands")
	chat.sendMsg("PASS oauth:" + t.config.Token.AccessToken)
	chat.sendMsg("NICK " + t.GetLoggedInUser().Login)

	go chat.readThread(conn)
}

func (c *Chat) sendMsg(message string) {
	_, err := c.Conn.Write([]byte(message + "\r\n"))
	if err != nil {
		log.Fatalf("Error sending message: %s", err)
	}
}

func (c *Chat) readThread(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		lineB, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("Error reading from server: %s", err)
		}
		line := string(lineB)

		if strings.Contains(line, ":tmi.twitch.tv 001 "+strings.ToLower(c.Twitch.GetLoggedInUser().Login)+" :Welcome, GLHF!") {
			// We've authenticated to the server
			log.Println("Connected!")
			c.Connected = true
			c.joinChannel()
		} else if strings.Contains(line, " 366 "+strings.ToLower(c.Twitch.GetLoggedInUser().Login)+" #"+strings.ToLower(c.Channel)) {
			// We've joined the desired channel
			log.Printf("Joined %s", c.Channel)
			c.Joined = true
		} else if strings.Contains(line, "PING") {
			// Respond to Keepalive message
			c.sendPong(line)
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG #"+strings.ToLower(c.Channel)+" :") {
			// Read a message in the streams chat
			m := parseMessage(line)
			fmt.Printf("%s: %s\n", m.Sender, m.Text)
		}
	}
}

func (c *Chat) joinChannel() {
	c.sendMsg("JOIN #" + c.Channel)
}

func (c *Chat) sendPong(line string) {
	rsp := strings.Replace(line, "PING", "PONG", 1)
	c.sendMsg(rsp)
}

func parseMessage(line string) *Message {
	m := new(Message)

	// Get the sender
	m.Sender = strings.Split(line, "!")[0][1:]

	// Get the message
	m.Text = strings.Split(line, "tmi.twitch.tv PRIVMSG #esfandtv :")[1]

	return m
}
