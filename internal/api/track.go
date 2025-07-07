package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

func GetTrack(ctx context.Context, client *spotify.Client, trackName string, trackID string, sleepDuration time.Duration) ([]spotify.FullTrack, error) {
	track, err := client.GetTrack(ctx, spotify.ID(trackName))
	if err != nil {
		return nil, fmt.Errorf("get track: %w", err)
	}

	log.Infof("acquired track `%s` (`%s`)", trackID, trackName)
	return []spotify.FullTrack{*track}, nil
}
