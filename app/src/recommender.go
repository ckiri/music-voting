package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)


type RecommendedSongs struct {
	songs []Song
}


func standardRequest(targetURL string) ([]byte, error) {

	request, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println("Failed creating a request:", err)
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed making request:", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed reading body:", err)
		return nil, err
	}

	fmt.Println("Request: ", targetURL, " Response: ", response.Status)
	return body, nil
}


func RecommendSongs(song Song, amount int) (RecommendedSongs, error) {

    queryURL := "https://ws.audioscrobbler.com/2.0/"
    method := "track.getsimilar"
    artistName := song.artistName
    songName := song.trackName
    apiKey := os.Getenv("LAST_FM_API_KEY")
    format := "json"
    limit := fmt.Sprintf("%d", amount)

    u, _ := url.Parse(queryURL)
	q := u.Query()
	q.Set("method", method)
	q.Set("artist", artistName)
	q.Set("track", songName)
	q.Set("api_key", apiKey)
	q.Set("format", format)
	q.Set("limit", limit)
	u.RawQuery = q.Encode()

	body, err := standardRequest(u.String())
	if err != nil {
		return RecommendedSongs{}, fmt.Errorf("Failed recommending song due to error: %s", err)
	}

	recommended := RecommendedSongs{}

	gjson.GetBytes(body, "similartracks.track").ForEach(func(_, value gjson.Result) bool {
        trackName := value.Get("name").String()
        artistName := value.Get("artist.name").String()
        searchString := fmt.Sprintf("%s %s", artistName, trackName)

        spotifySong, err := SearchForSong(searchString)
        if err != nil {
        	fmt.Println(err)
        	return true
        }
        recommended.songs = append(recommended.songs, spotifySong)
        return true
    })

    return recommended, nil

}

func AIRecommendSongs(song Song, amount int) (RecommendedSongs, error) {

    ollamaIP := os.Getenv("OLLAMA_IP")
    model := os.Getenv("OLLAMA_MODEL")
    limit := fmt.Sprintf("%d", amount)
    prompt := fmt.Sprintf(`You are a music recommendation system.
    Given the song "%s" by "%s", recommend %s similar songs.
    Only recommend more popular songs from the same genre and avoid little known artists.
    Return the result strictly as JSON with this format:
    {
      "recommended": [
        {"trackName": "...", "artistName": "..."},
        ...
      ]
    }`, song.trackName, song.artistName, limit)

    payload := map[string]any{
        "model": model,
        "prompt": prompt,
        "max_tokens": 10000,
    }
    payloadBytes, _ := json.Marshal(payload)

    ollamaURL := fmt.Sprintf("%s/v1/completions", ollamaIP)

    request, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
    	return RecommendedSongs{}, fmt.Errorf("Error creating request for Ollama: %s", err)
    }
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return RecommendedSongs{}, fmt.Errorf("Failed calling Ollama: %s", err)
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
    	return RecommendedSongs{}, fmt.Errorf("Failed reading body of Ollama response: %s", err)
    }

    llmOutput := gjson.GetBytes(body, "choices.0.text").String()

    recommended := RecommendedSongs{}

    gjson.Get(llmOutput, "recommended").ForEach(func(_, value gjson.Result) bool {
        trackName := value.Get("trackName").String()
        artistName := value.Get("artistName").String()
        searchString := fmt.Sprintf("%s %s", artistName, trackName)

        spotifySong, err := SearchForSong(searchString)
        if err != nil {
        	fmt.Println(err)
        	return true
        }
        recommended.songs = append(recommended.songs, spotifySong)
        return true
    })

    return recommended, nil
}
