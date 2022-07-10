package twitchgo

import (
	"regexp"
	"testing"
)

func TestParseConfig(t *testing.T) {
	configJSON := "{\"client_id\": \"MyID\", \"client_secret\": \"MySecret\"}"
	c, err := ParseConfig(configJSON)

	wantID := regexp.MustCompile("MyID")
	wantSecret := regexp.MustCompile("MySecret")

	if err != nil {
		t.Fatalf(`ParseConfig(configJSON) = %q, got error: %s`, c, err)
	}
	if !wantID.MatchString(c.ClientID) {
		t.Fatalf(`ParseConfig(configJSON) = got %s, want %s`, c.ClientID, wantID)
	} else if !wantSecret.MatchString(c.ClientSecret) {
		t.Fatalf(`ParseConfig(configJSON) = got %s, want %s`, c.ClientSecret, wantSecret)
	}
}
