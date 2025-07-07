package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

func GetPlaylist(ctx context.Context, client *spotify.Client, playlistName string, playlistID string, sleepDuration time.Duration) ([]spotify.FullTrack, error) {
	acquiredTracks := []spotify.FullTrack{}

	for {
		page, err := client.GetPlaylistItems(ctx, spotify.ID(playlistID))
		if err != nil {
			return nil, fmt.Errorf("get playlist items: %w", err)
		}

		for _, item := range page.Items {
			acquiredTracks = append(acquiredTracks, *item.Track.Track)
		}

		log.Infof("acquired %d tracks (tracks in page: %d) for playlist `%s` (`%s`)", len(acquiredTracks), len(page.Items), playlistID, playlistName)

		if page.Next == "" {
			break
		}
	}

	return acquiredTracks, nil
}
