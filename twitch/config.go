package twitch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Configuration struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Auth         string `json:"auth"`
	Token        Token  `json:"token"`
	path         string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	Expires      time.Time
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func LoadConfig(configPath string) *Configuration {
	// Read config file
	configJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the config JSON
	c, _ := ParseConfig(string(configJSON))
	c.path = configPath
	return c
}

func ParseConfig(configJSON string) (*Configuration, error) {
	c := new(Configuration)
	err := json.Unmarshal([]byte(configJSON), &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Configuration) WriteConfig(path string) {
	j, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Error marshaling config: %s", err)
	}

	log.Printf("%s\n", j)

	// Open config
	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalf("Error opening config at %s: %s", c.path, err)
	}

	// Write config
	_, err = f.Write(j)
	if err != nil {
		log.Fatalf("Error writing config: %s", err)
	}
	f.Close()
}
