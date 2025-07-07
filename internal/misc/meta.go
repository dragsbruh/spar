package misc

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
)

func SaveTracks(tracks []spotify.FullTrack, path string) error {
	bytes, err := json.MarshalIndent(tracks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(path, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func LoadTracks(path string) ([]spotify.FullTrack, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var tracks []spotify.FullTrack
	if err := json.Unmarshal(bytes, &tracks); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return tracks, nil
}
