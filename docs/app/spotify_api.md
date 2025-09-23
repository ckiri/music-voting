# Spotify Api + Recommender Usage

## Prerequisites

1. For inserting songs into queue - Spotify Premium:\
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

On start of the application, the function ```InitSpotify()``` **must** be called. This function ensures that the ```token.json``` file contains a valid bearer token. If this is the first time the app is started, you will need to copy the content in the stdout into your browser and authenticate to Spotify.

### Bugs / TODOs:
  - App crashes (sometimes ?) after the first auth
  - Sometimes the refresh works but the function thereafter still uses the old key - after a restart the code works again
  - Port of callback can overlap with the gin web deployment
  - The import of the .env file is relative and this should not be

## Spotify API (spotify_api.go)

Once ```InitSpotify()``` has been called, the following functions are given to the user:

```go
func GetPlaybackState() PlaybackState { ... }
```
Returns a PlaybackState struct witch contains the following information:
```go
type PlaybackState struct {
	song 			Song		//Song struct below
	playbackState 	string		// (true|false) is the song playing?
	progress 		int64		//leftover song time in ms
}
// Playback contains struct song:
type Song struct {
	trackName		string	// Name of the track
	trackId			string	// SpotifyID of the track
	artistName 	string	// Artist name
	artistId 		string	// SpotifyID of the artist
	uri				string	// uri used when telling spotify what song to play
}
```
Also the following (self-explanatory) functions can be called:
```go
func AddSongToQueue(songURI string) { ... }
```
```go
func SkipCurrentSong() { ... }
```
```go
func SearchForSong(searchString string) Song { ... }
```

### Bugs / TODOs:
- No error handling whatsoever -> unexpected behaviour when a API call inevitably fails
- Could siplify the function parameters to always accept Song structs and handle value extraction in code
- More information could be collected, like album, other artists ...

## Recommending functions (recommender.go)

In order to handle the logic of recommending a song based on the actual playing song the module ```recommender.go``` is written. Due to spotify depricating the ```/recomendations``` [:'(](https://developer.spotify.com/documentation/web-api/reference/get-recommendations) a workaround is made using the information of the playing song and the [Last.fm API](https://www.last.fm/api).

The following function is available to call:

```go
func RecommendSongs(song Song) RecommendedSongs { ...}
```

which returns a struct of possible songs to play (5 in total):

```go
type RecommendedSongs struct {
	Songs []Song
}
```
Alternatively the following ✨**AI Powered**✨ function is made available:

```go
func AIRecommendSongs(song Song) RecommendedSongs { ... }
```

Same logic, just worse (performace / latency / hallucinations). You also need a local OLLAMA Instance running as well as a comprehensive model that does not hallucinate songs into existance.

### Bugs / TODOs:
- No error handling whatsoever -> unexpected behaviour when a API call inevitably fails
- More choices for caller (e.g. how many songs to recommend, how big of a token context ...)

## Docs

[Spotify API for Developers](https://developer.spotify.com)\
[LastFM API for getting similar songs](https://www.last.fm/api/show/track.getSimilar)\
[Ollama for LLM powered recommendations](https://ollama.com)
