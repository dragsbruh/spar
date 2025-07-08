package api

import (
	"context"
	"fmt"
	"time"

	"github.com/dragsbruh/spar/internal/listfile"
	"github.com/zmb3/spotify/v2"
)

func GetItem(ctx context.Context, client *spotify.Client, item listfile.ListfileItem, sleepDuration time.Duration) ([]spotify.FullTrack, error) {
	switch item.Kind {
	case "artist":
		return GetArtist(ctx, client, item.Name, item.ID, sleepDuration)
	case "playlist":
		return GetPlaylist(ctx, client, item.Name, item.ID, sleepDuration)
	case "track":
		return GetTrack(ctx, client, item.Name, item.ID, sleepDuration)
	default:
	}

	return nil, fmt.Errorf("unreachable")
}
