package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)


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


func makeRequest(targetURL string, requestType string) (string, error) {

	token, err := ReadToken(TokenFile)

	request, err := http.NewRequest(requestType, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("Failed creating a request: %s", err)
	}
	request.Header.Add("Authorization", "Bearer " + token.AccessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("Failed making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Failed reading body: %s", err)
	}

	fmt.Println("Status: ", response.Status)
	return string(body), nil
}


func GetPlaybackState() (PlaybackState, error) {

	body, err := makeRequest("https://api.spotify.com/v1/me/player", "GET")
	if err != nil {
		return PlaybackState{}, err
	}

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

	return playbackState, nil
}


func AddSongToQueue(song Song) error {

	songURI := song.uri
 	targetURL := fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", songURI)
	_, err := makeRequest(targetURL, "POST")

	if err != nil {
		return err
	} else {
		return nil
	}
}


func SkipCurrentSong() error {

	_, err := makeRequest("https://api.spotify.com/v1/me/player/next", "POST")
	if err != nil {
		return err
	} else {
		return nil
	}
}


func SearchForSong(searchString string) (Song, error) {

	searchString = url.QueryEscape(searchString)
	targetURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1", searchString)
	body, err := makeRequest(targetURL, "GET")
	if err != nil {
		return Song{}, err
	}

	song := Song {
		trackName: 		gjson.Get(body, "tracks.items.0.name").String(),
		trackId: 		gjson.Get(body, "tracks.items.0.id").String(),
		artistName: 	gjson.Get(body, "tracks.items.0.artists.0.name").String(),
		artistId: 		gjson.Get(body, "tracks.items.0.artists.0.id").String(),
		uri:			gjson.Get(body, "tracks.items.0.uri").String(),
	}

	return song, nil
}
