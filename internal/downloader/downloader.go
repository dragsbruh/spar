package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

func DownloadTracks(tracks []spotify.FullTrack, tempDir, outDir string, maxWorkers int) error {
	logFile, err := os.Create(filepath.Join(tempDir, ".spar.thirdparty.log"))
	if err != nil {
		return fmt.Errorf("creating log file: %w", err)
	}
	defer logFile.Close()

	type result struct {
		index int
		err   error
	}

	jobs := make(chan int, len(tracks))
	results := make(chan result, len(tracks))

	for range maxWorkers {
		go func() {
			for i := range jobs {
				err := DownloadSingleTrack(tracks[i], tempDir, outDir, logFile)
				results <- result{i, err}
			}
		}()
	}

	for i := range tracks {
		jobs <- i
	}
	close(jobs)

	for range tracks {
		res := <-results
		if res.err != nil {
			log.Errorf("(%d/%d) failed to download `%s`: %v", res.index+1, len(tracks), tracks[res.index].Name, res.err)
		} else {
			log.Infof("(%d/%d) downloaded audio for `%s` (`%s`)", res.index+1, len(tracks), tracks[res.index].ID, tracks[res.index].Name)
		}
	}

	return nil
}

func DownloadSingleTrack(track spotify.FullTrack, tempDir string, outDir string, logFile io.Writer) error {
	rawAudioPath := filepath.Join(tempDir, fmt.Sprintf("raw_%s.mp3", track.ID))
	rawCoverPath := filepath.Join(tempDir, fmt.Sprintf("cover_%s.jpg", track.ID))
	finalAudioPath := filepath.Join(outDir, fmt.Sprintf("%s - %s.mp3", slug.Make(track.Artists[0].Name), slug.Make(track.Name)))

	_, err := os.Stat(finalAudioPath)
	if err == nil {
		return nil
	}

	hasCover := len(track.Album.Images) > 0

	if hasCover {
		file, err := os.Create(rawCoverPath)
		if err != nil {
			return fmt.Errorf("opening cover path for %s: %w", track.ID, err)
		}
		defer file.Close()

		if err := track.Album.Images[0].Download(file); err != nil {
			return fmt.Errorf("downloading cover for %s: %w", track.ID, err)
		}
	}

	if err := DownloadAudio(track, rawAudioPath, logFile); err != nil {
		return fmt.Errorf("downloading audio for %s: %w", track.ID, err)
	}

	if err := AddMetadata(track, rawAudioPath, rawCoverPath, finalAudioPath, logFile); err != nil {
		return fmt.Errorf("adding metadata for %s: %w", track.ID, err)
	}

	return nil
}
