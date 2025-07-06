package tokenutil

import (
	"encoding/json"
	"os"

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

func Load(path string) *oauth2.Token {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	defer file.Close()

	var token oauth2.Token
	json.NewDecoder(file).Decode(&token)
	return &token
}
