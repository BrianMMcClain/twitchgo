package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Configuration struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthToken    string `json:"auth_token"`
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
