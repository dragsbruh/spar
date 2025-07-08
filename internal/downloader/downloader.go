package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dragsbruh/spar/internal/misc"
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
	audioFilePathWithoutExt := filepath.Join(tempDir, fmt.Sprintf("audio_%s", track.ID))
	rawMetadataPath := filepath.Join(tempDir, fmt.Sprintf("meta_%s.json", track.ID))
	symlinkCoverPath := filepath.Join(tempDir, fmt.Sprintf("cover_%s.jpg", track.ID))
	albumCoverPath := filepath.Join(tempDir, fmt.Sprintf("acover_%s.jpg", track.Album.ID))

	rawAudioPath, err := misc.GetAudioFiles(audioFilePathWithoutExt)
	if err != nil {
		return fmt.Errorf("getting audio files: %w", err)
	}

	audioExists := rawAudioPath != ""

	_, err = os.Stat(symlinkCoverPath)
	coverExists := err == nil

	_, err = os.Stat(rawMetadataPath)
	metaExists := err == nil

	_, err = os.Stat(albumCoverPath)
	albumCoverExists := err == nil

	hasCover := len(track.Album.Images) > 0

	if hasCover {
		if !albumCoverExists {
			file, err := os.Create(albumCoverPath)
			if err != nil {
				return fmt.Errorf("opening cover path for %s: %w", track.Album.ID, err)
			}
			defer file.Close()

			if err := track.Album.Images[0].Download(file); err != nil {
				return fmt.Errorf("downloading cover for %s: %w", track.ID, err)
			}
		}
		if !coverExists {
			if err := os.Symlink(albumCoverPath, symlinkCoverPath); err != nil {
				return fmt.Errorf("creating symlink cover for %s: %w", track.ID, err)
			}
		}
	}

	if !audioExists {
		if err := DownloadAudio(track, fmt.Sprintf("%s.%%(ext)s", audioFilePathWithoutExt), logFile); err != nil {
			return fmt.Errorf("downloading audio for %s: %w", track.ID, err)
		}
	}

	if !metaExists {
		if err := SaveMetadata(track, rawMetadataPath); err != nil {
			return fmt.Errorf("adding metadata for %s: %w", track.ID, err)
		}
	}

	return nil
}
