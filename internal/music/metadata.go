package music

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

func GetItemTracks(ctx context.Context, client *spotify.Client, item ListItem, sleepDuration time.Duration) []spotify.FullTrack {
	switch item.Kind {
	case ListItemTrack:
		defer time.Sleep(sleepDuration)
		track, err := client.GetTrack(ctx, spotify.ID(item.Id))
		if err != nil {
			log.Warnf("Failed to get track item `%s` (id: `%s`): %v", item.Name, item.Id, err)
			return nil
		}
		return []spotify.FullTrack{*track}

	case ListItemPlaylist:
		acquiredTracks := []spotify.FullTrack{}

		for {
			time.Sleep(sleepDuration)
			page, err := client.GetPlaylistItems(ctx, spotify.ID(item.Id), spotify.Offset(len(acquiredTracks)))
			if err != nil {
				log.Warnf("Failed to get playlist page for `%s` (id: `%s`, acquired tracks %d before error): %v", item.Name, item.Id, len(acquiredTracks), err)
				continue // TODO: this might result in infinite loop if we always error
			}
			for _, item := range page.Items {
				acquiredTracks = append(acquiredTracks, *item.Track.Track)
			}

			log.Infof("Acquired %d tracks for page of playlist `%s` (id: `%s`, total acquired: %d)", len(page.Items), item.Name, item.Id, len(acquiredTracks))
			if page.Next == "" {
				break
			}
		}

		return acquiredTracks

	case ListItemArtist:
		acquiredAlbums := []spotify.SimpleAlbum{}

		for {
			time.Sleep(sleepDuration)
			albums, err := client.GetArtistAlbums(ctx, spotify.ID(item.Id), []spotify.AlbumType{spotify.AlbumTypeAlbum}, spotify.Offset(len(acquiredAlbums)))
			if err != nil {
				log.Warnf("Failed to get artist albums for `%s` (id: `%s`): %v", item.Name, item.Id, err)
				continue // TODO: this too
			}
			log.Infof("Acquired %d albums for page of artist `%s` (id: `%s`, total acquired albums: %d)", len(albums.Albums), item.Name, item.Id, len(acquiredAlbums))
			acquiredAlbums = append(acquiredAlbums, albums.Albums...)
			if albums.Next == "" {
				break
			}
		}

		acquiredTracks := []spotify.FullTrack{}

		for i, album := range acquiredAlbums {
			time.Sleep(sleepDuration)
			album, err := client.GetAlbum(ctx, album.ID)
			if err != nil {
				log.Warnf("Failed to get album %s for `%s` (id: `%s`): %v", album.Name, item.Name, item.Id, err)
				continue
			}

			for _, track := range album.Tracks.Tracks {
				acquiredTracks = append(acquiredTracks, spotify.FullTrack{
					SimpleTrack: spotify.SimpleTrack{
						Artists:          track.Artists,
						AvailableMarkets: track.AvailableMarkets,
						DiscNumber:       track.DiscNumber,
						Duration:         track.Duration,
						Explicit:         track.Explicit,
						ExternalURLs:     track.ExternalURLs,

						Endpoint:    track.Endpoint,
						ID:          track.ID,
						Name:        track.Name,
						PreviewURL:  track.PreviewURL,
						TrackNumber: track.TrackNumber,
						URI:         track.URI,
						Type:        track.Type,
					},
					Album:       album.SimpleAlbum,
					ExternalIDs: map[string]string{}, // meh who needs this anyway
				})
			}

			log.Infof("Acquired %d tracks for album `%s` (%d/%d) for artist `%s` (id: `%s`, total acquired: %d)", len(album.Tracks.Tracks), album.Name, i+1, len(acquiredAlbums), item.Id, item.Id, len(acquiredTracks))
		}

		return acquiredTracks

	default:
		log.Fatal("Unknown item type")
		return nil // unreachable
	}
}
