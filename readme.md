# Spar

> **NOTE:** this should be functional but is still in development

Spar stands for **Sp**otify **Ar**chiver.
Written in golang, its supposed to be archiving playlists/tracks/artists.

Unlike spotdl, spar uses Spotify API to get the metadata and yt-dlp to download the actual audio from YouTube.

I suck at error handling, beware of errors, Although they should be rare.

## How it works

You define the artists and stuff declaratively in `spar.yml` like this:

**Format:**

```yaml
tempdir: ~/temp/music_cache/ # temporary place to store cover files, raw audio etc
outdir: ~/Music/ # directory in which music is downloaded to as mkv
metapath: ~/Music/.meta.json # metadata json file (useful if you want to pause downloading and restart again, use with `-l` flag after you run it once)
workers: 10 # number of concurrent downloads

items:
  - name: My fav artist # artist name, not used internally but displayed
    kind: artist # available: artist/playlist/track
    id: <artist_id> # just the id, ex: `2r3DAIz6afSzxVnM1Rzj3N`

  - name: My Playlist
    kind: playlist
    id: <playlist_id>  # note that curated playlists and stuff might not work
```

### Downloading

It downloads tracks parallely using yt-dlp by searching for the closest match to the spotify track.

## Installation

You must have `go` installed, and a Spotify Account (for the API).

1. Clone this repository.
2. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard) and create a new app.
3. Fill in your preferred details and add this redirect URI: `http://127.0.0.1:8080/callback` (very important).
4. After creation, copy your client ID and client secret and save them in your `.env` (or environment) as `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET`.
5. Run `go build ./cmd/main/main.go`

The first time you run `spar` it will prompt you to authenticate to Spotify via OAuth, and then save the token as `.spar-token` in your user `.config/spar` directory.
You will need to reauthenticate if that file is missing or token expired.

## Usage

In the same directory as your `spar.yml`

```bash
spar sync
```

Thats it

**Continuing downloads:**

In most cases, you can just ctrl+c if you want to go somewhere.
When you are back, just hit

```bash
spar sync --local
```

This will load the tracks from your local metadata json instead of hitting the Spotify API again.

**Using custom listfile:**

```bash
spar sync --listfile=mylist.yaml
```

or `--file` or `--lf` flags work too.

## TO-DO

*(In order of priority)*

- [x] Download the actual music (to a folder with the sluggified artist-title as filename)
- [ ] CLI config for music format, list file, rate limit handling, etc
- [ ] Make it feel like a legit app by using better logs
- [ ] Use as a standlone app (directly specify)

## Docs

**Flow:**

1. Get tracks from Spotify API.
2. Iterate through tracks and save covers to temp directory, get the raw opus audio from youtube and embed metadata with audio in an mkv container.
