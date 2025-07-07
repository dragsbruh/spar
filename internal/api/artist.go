package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

func GetArtist(ctx context.Context, client *spotify.Client, artistName string, artistID string, sleepDuration time.Duration) ([]spotify.FullTrack, error) {
	acquiredAlbums := []spotify.SimpleAlbum{}

	for {
		page, err := client.GetArtistAlbums(ctx, spotify.ID(artistID), []spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle}, spotify.Offset(len(acquiredAlbums)))
		if err != nil {
			return nil, fmt.Errorf("getting artist albums (%s): %w", artistID, err)
		}

		acquiredAlbums = append(acquiredAlbums, page.Albums...)

		log.Infof("acquired %d albums (albums in page: %d) for artist `%s` (`%s`)", len(acquiredAlbums), len(page.Albums), artistID, artistName)

		if page.Next == "" {
			break
		}

		time.Sleep(sleepDuration)
	}

	acquiredTracks := []spotify.FullTrack{}

	for i, album := range acquiredAlbums {
		acquiredAlbumTracks := []spotify.FullTrack{}

		for {
			page, err := client.GetAlbumTracks(ctx, album.ID, spotify.Offset(len(acquiredAlbumTracks)))
			if err != nil {
				return nil, fmt.Errorf("get album tracks (%s -> %s): %w", artistID, album.ID, err)
			}

			for _, track := range page.Tracks {
				ftrack := spotify.FullTrack{
					SimpleTrack: track,
					Album:       album,
				}
				track.Album = album
				acquiredAlbumTracks = append(acquiredAlbumTracks, ftrack)
			}

			log.Infof("acquired %d tracks (tracks in page: %d) for album `%s` for artist `%s` (`%s`, `%s`)", len(acquiredAlbumTracks), len(page.Tracks), album.ID, artistID, album.Name, artistName)

			time.Sleep(sleepDuration)

			if page.Next == "" {
				break
			}

		}

		log.Infof("acquired %d tracks (tracks in album: %d, album %d/%d) for artist `%s` (`%s`)", len(acquiredTracks), len(acquiredAlbumTracks), i+1, len(acquiredAlbums), artistID, artistName)

		acquiredTracks = append(acquiredTracks, acquiredAlbumTracks...)
	}

	return acquiredTracks, nil
}
