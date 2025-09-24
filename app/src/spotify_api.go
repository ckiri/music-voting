package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

var (
	token, err = ReadToken(TokenFile)
)

func makeRequest(targetURL string, requestType string) string {
	request, err := http.NewRequest(requestType, targetURL, nil)
	if err != nil {
		fmt.Println("Failed creating a request:", err)
		return ""
	}
	request.Header.Add("Authorization", "Bearer " + token.AccessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed making request:", err)
		return ""
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed reading body:", err)
		return ""
	}

	fmt.Println("Status:", response.Status)
	return string(body)
}


type Song struct {
	trackName	string
	trackId		string
	artistName 	string
	artistId 	string
	uri			string
}


type PlaybackState struct {
	song 			Song
	playbackState 	string
	progress 		int64
}


func GetPlaybackState() PlaybackState {

	body := makeRequest("https://api.spotify.com/v1/me/player", "GET")

	progressMs := gjson.Get(body, "progress_ms").Int()
	durationMs := gjson.Get(body, "item.duration_ms").Int()
	timeLeftMs := durationMs - progressMs

	song := Song {
		trackName: 		gjson.Get(body, "item.name").String(),
		trackId: 		gjson.Get(body, "item.id").String(),
		artistName: 	gjson.Get(body, "item.artists.0.name").String(),
		artistId: 		gjson.Get(body, "item.artists.0.id").String(),
		uri:			gjson.Get(body, "item.uri").String(),
	}

	playbackState := PlaybackState{
		song:			song,
		playbackState: 	gjson.Get(body, "is_playing").String(),
		progress: 		timeLeftMs,
	}

	return playbackState
}

func AddSongToQueue(song Song) {

	songURI := song.uri
 	targetURL := fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", songURI)
	makeRequest(targetURL, "POST")
}

func SkipCurrentSong() {

	makeRequest("https://api.spotify.com/v1/me/player/next", "POST")
}

func SearchForSong(searchString string) Song {

	searchString = url.QueryEscape(searchString)
	targetURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1", searchString)
	body := makeRequest(targetURL, "GET")

	song := Song {
		trackName: 		gjson.Get(body, "tracks.items.0.name").String(),
		trackId: 		gjson.Get(body, "tracks.items.0.id").String(),
		artistName: 	gjson.Get(body, "tracks.items.0.artists.0.name").String(),
		artistId: 		gjson.Get(body, "tracks.items.0.artists.0.id").String(),
		uri:			gjson.Get(body, "tracks.items.0.uri").String(),
	}

	return song
}
