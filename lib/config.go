package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type Configuration struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Auth         AuthCode `json:"auth"`
	Token        Token    `json:"token"`
}

type AuthCode struct {
	auth_code string `json:"auth_code"`
	token     string `json:"auth_token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	Expires      time.Time
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func LoadConfig(configPath string) *Configuration {
	c := new(Configuration)

	// Read config file
	configJson, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the config JSON
	json.Unmarshal([]byte(configJson), &c)
	return c
}

func (c *Configuration) WriteConfig(path string) {
	j, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Error writing config: %s", err)
	}
	fmt.Println(string(j))
}
