package downloader

import (
	"fmt"
	"os/exec"

	"github.com/zmb3/spotify/v2"
)

func DownloadOpusAudio(track spotify.FullTrack, output string) error {
	args := []string{
		"-f", "bestaudio",
		"--extract-audio",
		"--audio-format", "opus",
		"--output", output,
		fmt.Sprintf("ytsearch1:%s %s", track.Artists[0].Name, track.Name),
	}

	cmd := exec.Command("yt-dlp", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp error: %w", err)
	}

	return nil
}

func AddMetadata(track spotify.FullTrack, rawAudioPath string, rawCoverPath string, outputPath string) error {
	metadataArgs := []string{
		"-metadata", fmt.Sprintf("title=%s", track.Name),
		"-metadata", fmt.Sprintf("artist=%s", track.Artists[0].Name),
		"-metadata", fmt.Sprintf("album=%s", track.Album.Name),
		"-metadata", fmt.Sprintf("date=%s", track.Album.ReleaseDate),
		"-metadata", fmt.Sprintf("track=%d", track.TrackNumber),
	}
	if len(track.Album.Artists) > 0 {
		metadataArgs = append(metadataArgs,
			"-metadata", fmt.Sprintf("album_artist=%s", track.Album.Artists[0].Name),
		)
	}

	args := []string{
		"-y",
		"-i", rawAudioPath,
		"-i", rawCoverPath,
		"-map", "0:a",
		"-map", "1",
		"-c:a", "libopus",
		"-b:a", "192k",
		"-c:v", "copy",
		"-metadata:s:v", "title=Album cover",
		"-metadata:s:v", "comment=Cover (front)",
	}
	args = append(args, metadataArgs...)
	args = append(args, outputPath)

	cmd := exec.Command("ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}

	return nil
}
