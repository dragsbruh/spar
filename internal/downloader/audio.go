package downloader

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/zmb3/spotify/v2"
)

func DownloadOpusAudio(track spotify.FullTrack, output string, logFile io.Writer) error {
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

func AddMetadata(track spotify.FullTrack, rawAudioPath string, rawCoverPath string, outputPath string, logFile io.Writer) error {
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
		"-map", "0",
		"-c", "copy",
		"-metadata:s:v", `title=Album cover`,
		"-metadata:s:v", `comment=Cover (front)`,
		"-disposition:v", "attached_pic",
	}

	args = append(args, metadataArgs...)
	args = append(args, outputPath)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}

	return nil
}
