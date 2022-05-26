package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewConfig(configPath string) *Configuration {
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
