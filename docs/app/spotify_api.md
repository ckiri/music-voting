# Spotify Api + Recommender Usage

## Prerequisites

1. For inserting songs into queue - Spotify Premium:
- Create new email
- Open private window in your browser
- Get a month of free spotify
- repeat
2. go + following packages
- [gjson](https://github.com/tidwall/gjson) - Easy json parsing
- [godotenv](https://github.com/joho/godotenv) - Loads environment variables from .env files
- [oauth2](https://golang.org/x/oauth2) - client implementation for OAuth 2.0 spec
3. .env file filled (use .env.example for reference)

## Initialising and Auth (spotify_auth.go)

On start of the application, the function ```InitSpotify()``` **must** be called. This function ensures that the ```token.json``` file contains a valid bearer token. If this is the first time the app is started, you will need to copy the content in the stdout into your browser and authenticate to Spotify. If any error happens in ```InitSpotify()``` which leads to the token not being valid, the program crashes.

### Bugs / TODOs
  - Port of callback can overlap with the gin web deployment
  - Environment variable handling (e.g. ```loadEnv()```) should be development only and not called always. In deployment these values can be setup in Docker and the code should only check if the values are present, not read a .env file

## Spotify API (spotify_api.go)

Once ```InitSpotify()``` has been called, the following functions are given to the user:

```go
func GetPlaybackState() (PlaybackState, error) { ... }
```
Returns a PlaybackState struct witch contains the following information:
```go
type PlaybackState struct {
	song            Song		// Song struct below
	playbackState   string		// (true|false) is the song playing?
	progress        int64		// leftover song time in ms
}
// Playback contains struct song:
type Song struct {
	trackName   string	// Name of the track
	trackId     string	// SpotifyID of the track
	artistName  string	// Artist name
	artistId    string	// SpotifyID of the artist
	uri         string	// uri used when telling spotify what song to play
}
```
Also the following (self-explanatory) functions can be called:
```go
func AddSongToQueue(song Song) error { ... }
```
```go
func SkipCurrentSong() error { ... }
```
```go
func SearchForSong(searchString string) (Song, error) { ... }
```

### Bugs / TODOs
- More information could be collected, like album, other artists ...

## Recommending functions (recommender.go)

In order to handle the logic of recommending a song based on the actual playing song the module ```recommender.go``` is written. Due to Spotify deprecating the ```/recommendations``` [:'(](https://developer.spotify.com/documentation/web-api/reference/get-recommendations) a workaround is made using the information of the playing song and the [Last.fm API](https://www.last.fm/api).

The following function is available to call:

```go
func RecommendSongs(song Song, amount int) (RecommendedSongs, error) { ... }
```

which returns a struct of possible songs to play (5 in total):

```go
type RecommendedSongs struct {
	Songs []Song
}
```
Alternatively the following ✨**AI Powered**✨ function is made available:

```go
func AIRecommendSongs(song Song, amount int) (RecommendedSongs, error) { ... }
```

Same logic, just worse (performance / latency / hallucinations). You also need a local OLLAMA Instance running as well as a comprehensive model that does not hallucinate songs into existence.

### Bugs / TODOs
- More choices for caller (e.g. context size, proompting, genre limits etc.)

## Error handling

All public facing functions witch **can** fail return a ```error```-type object. This should be checked after calling the function. The error object can be printed out to get more information like a stacktrace.

## Docs

[Spotify API for Developers](https://developer.spotify.com)\
[LastFM API for getting similar songs](https://www.last.fm/api/show/track.getSimilar)\
[Ollama for LLM powered recommendations](https://ollama.com)
