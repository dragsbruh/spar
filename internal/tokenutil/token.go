package tokenutil

import (
	"encoding/json"
	"os"
	"time"

	"golang.org/x/oauth2"
)

const TokenFile = ".spar_token"

func Save(path string, token *oauth2.Token) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	json.NewEncoder(file).Encode(token)
	file.Close()
	return nil
}

// returns token, needsRefresh
func Load(path string) (*oauth2.Token, bool) {
	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}

	defer file.Close()

	var token oauth2.Token
	json.NewDecoder(file).Decode(&token)

	if time.Now().After(token.Expiry) {
		return &token, true
	}

	return &token, false
}
