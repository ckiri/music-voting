package main

import (
	"fmt"
	"os"
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"bytes"

	"github.com/tidwall/gjson"
)

func standardRequest(targetURL string) []byte {

	request, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println("Failed creating a request:", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed making request:", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed reading body:", err)
	}

	fmt.Println("Status:", response.Status)
	return body
}

type RecommendedSongs struct {
	Songs []Song
}


func RecommendSongs(song Song, amount int) RecommendedSongs {
	/*
	* Gets a song and artist running rn
	* Searches lastfm with the name and gets 5 songs back
	* Searches for these songs on spotify
	* TODO: Abstract the query
    */

    songName := song.trackName
    artistName := song.artistName
    apiKey := os.Getenv("LAST_FM_API_KEY")
    queryURL := "https://ws.audioscrobbler.com/2.0/"
    limit := fmt.Sprintf("%d", amount)

    u, _ := url.Parse(queryURL)
	q := u.Query()
	q.Set("method", "track.getsimilar")
	q.Set("artist", artistName)
	q.Set("track", songName)
	q.Set("api_key", apiKey)
	q.Set("format", "json")
	q.Set("limit", limit)
	u.RawQuery = q.Encode()

	recommended := RecommendedSongs{}

	body := standardRequest(u.String())

	gjson.GetBytes(body, "similartracks.track").ForEach(func(_, value gjson.Result) bool {
        trackName := value.Get("name").String()
        artistName := value.Get("artist.name").String()
        searchString := fmt.Sprintf("%s %s", artistName, trackName)

        spotifySong := SearchForSong(searchString)
        recommended.Songs = append(recommended.Songs, spotifySong)
        return true
    })

    return recommended

}

func AIRecommendSongs(song Song, amount int) RecommendedSongs {
	/*
	* Gets a song and artist running rn
	* Asks an LLM for 5 song / artist combos
	* Searches for these songs on spotify
    */

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
    fmt.Println(prompt)
    payload := map[string]any{
        "model": model,
        "prompt": prompt,
        "max_tokens": 10000,
    }
    payloadBytes, _ := json.Marshal(payload)
    url := fmt.Sprintf("%s/v1/completions", ollamaIP)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("failed calling Ollama: %v", err)
        return RecommendedSongs{}
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    fmt.Println(string(body))

    llmOutput := gjson.GetBytes(body, "choices.0.text").String()
    fmt.Println(llmOutput)

    recommended := RecommendedSongs{}

    gjson.Get(llmOutput, "recommended").ForEach(func(_, value gjson.Result) bool {
        trackName := value.Get("trackName").String()
        artistName := value.Get("artistName").String()
        searchString := fmt.Sprintf("%s %s", artistName, trackName)

        spotifySong := SearchForSong(searchString)
        recommended.Songs = append(recommended.Songs, spotifySong)
        return true
    })

    return recommended
}
