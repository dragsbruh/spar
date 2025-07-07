package tokenutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

func TokenPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homedir, ".config", "spar", "token.json"), nil
}

func Load() (*oauth2.Token, error) {
	tokenPath, err := TokenPath()
	if err != nil {
		return nil, fmt.Errorf("getting token path: %w", err)
	}

	source, err := os.ReadFile(tokenPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading token file: %w", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(source, &token); err != nil {
		return nil, fmt.Errorf("unmarshalling token json: %w", err)
	}

	return &token, nil
}

func Save(token *oauth2.Token) error {
	tokenPath, err := TokenPath()
	if err != nil {
		return fmt.Errorf("getting token path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(tokenPath), os.ModePerm); err != nil {
		return fmt.Errorf("creating directory for token: %w", err)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshalling token: %w", err)
	}

	if err := os.WriteFile(tokenPath, data, os.ModePerm); err != nil {
		return fmt.Errorf("writing token json: %w", err)
	}

	return nil
}
