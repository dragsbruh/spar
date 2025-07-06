# Spar

> **NOTE:** not functional yet, only gets full metadata as of now.

Spar stands for **Sp**otify **Ar**chiver.
Written in golang, its supposed to be archiving playlists/tracks/artists.

Unlike spotdl, spar uses Spotify API to get the metadata and yt-dlp to download the actual audio from YouTube.

I suck at error handling, beware of errors, Although they should be rare.

## How it works

You define the artists and stuff in `list.csv` like this:

**Format:**

```csv
Name,Type,ID
```

**Fields:**

1. **Name:** Friendly name, not used internally
2. **Type:** Type of item, can be `playlist`, `track` or `artist` (case insensitive)
3. **ID:** Spotify ID of item, can be found at the end of path of URL of playlist/track/artist

**Example:**

```csv
HOYO-MiX,artist,2YvlK6lKiKVjXxsjvNbnqg
SawanoHiroyuki,artist,0Riv2KnFcLZA3JSVryRg4y
SawanoHiroyuki [nZk],artist,2EWXgN0xWOnbqJOxa9pWNO
Krage,artist,35jRIUtWCUITFLfjhYwkFx
My Beloved,playlist,6PQh5zfIng0G3bDkHmvjK7
```

Please do not use comments I did not implement that.

### Downloading

Similar to spotdl, it 

It downloads tracks parallely using yt-dlp

## Installation

You must have `go` and `git` installed, and a Spotify Account (for the API).

1. Clone this repository.
2. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard) and create a new app.
3. Fill in your preferred details and add this redirect URI: `http://127.0.0.1:8080/callback` (very important).
4. After creation, copy your client ID and client secret and save them in your `.env` (or environment) as `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET`.
5. Run `go build ./cmd/main/main.go`

## Usage

In the same directory as your `list.csv`

```bash
spar
```

Thats it

**NOTE:** Actual UX, Configuration is To-Do

## TO-DO

*(In order of priority)*

- [ ] Download the actual music (to a folder with the track id as filename)
- [ ] CLI config for music format, list file, rate limit handling, etc
- [ ] Make it feel like a legit app by using better logs
- [ ] Use as a standlone app (directly specify)
