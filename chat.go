package twitchgo

import (
	"bufio"
	"log"
	"net"
	"strconv"
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
	Sender     string
	Text       string
	Subscriber bool
	SubLength  int
	Mod        bool
	UserID     string
	Channel    string
}

var chatCallback func(*Message)

func (t *Twitch) ChatConnect(channel string, handler func(*Message)) {
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
	chat.sendMsg("CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands")
	chat.sendMsg("PASS oauth:" + t.config.Token.AccessToken)
	chat.sendMsg("NICK " + t.GetLoggedInUser().Login)

	chatCallback = handler

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
			c.Connected = true
			c.joinChannel()
		} else if strings.Contains(line, " 366 "+strings.ToLower(c.Twitch.GetLoggedInUser().Login)+" #"+strings.ToLower(c.Channel)) {
			// We've joined the desired channel
			c.Joined = true
		} else if strings.Contains(line, "PING") {
			// Respond to Keepalive message
			c.sendPong(line)
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG #"+strings.ToLower(c.Channel)+" :") {
			// Read a message in the streams chat
			m := c.parseMessage(line)
			chatCallback(m)
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

func (c *Chat) parseMessage(line string) *Message {
	m := new(Message)
	m.Channel = c.Channel

	// Parse the advanced tags to pull the user, sub, and mod info
	tags := strings.Split(strings.Split(line, "!")[0][1:], ";")
	for _, tag := range tags {
		tagSplit := strings.Split(tag, "=")
		key := tagSplit[0]
		val := tagSplit[1]
		if key == "display-name" {
			m.Sender = val
		} else if key == "subscriber" {
			m.Subscriber = val == "1"
		} else if key == "badge-info" {
			// Sub length only exists if they are a sub, ensure it exists
			if len(val) > 0 {
				if strings.Split(val, "/")[0] == "subscriber" {
					subLength, err := strconv.Atoi(strings.Split(val, "/")[1])
					if err != nil {
						log.Fatalf("Could not determine sub length: %s\n", err)
					}
					m.SubLength = subLength
				}
			}
		} else if key == "mod" {
			m.Mod = val == "1"
		} else if key == "user-id" {
			m.UserID = val
		}
	}

	// Get the message
	m.Text = strings.Split(line, "tmi.twitch.tv PRIVMSG #"+c.Channel+" :")[1]

	return m
}
