package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/zmb3/spotify/v2"
)

func DownloadAudio(track spotify.FullTrack, output string, logFile io.Writer) error {
	args := []string{
		"-f", "bestaudio",
		"--no-playlist",
		"--no-mtime",
		"--output", output,
		fmt.Sprintf("ytsearch1:%s %s", track.Artists[0].Name, track.Name),
	}

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp error: %w", err)
	}

	return nil
}

func SaveMetadata(track spotify.FullTrack, metadataPath string) error {
	metadata := map[string]string{
		"title":  track.Name,
		"artist": track.Artists[0].Name,
		"album":  track.Album.Name,
		"date":   track.Album.ReleaseDate,
		"track":  fmt.Sprintf("%d", track.TrackNumber),
	}

	bytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshalling json: %w", err)
	}

	if err := os.WriteFile(metadataPath, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("writing file %s: %w", metadataPath, err)
	}

	return nil
}
