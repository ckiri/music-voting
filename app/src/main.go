package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	demo()
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}

func demo() {

	InitSpotify()
	playbackState := GetPlaybackState()
	fmt.Println("Player state: ", playbackState)

	song := playbackState.song

	recommendedSongs := RecommendSongs(song, 5)
	fmt.Println("Recommended songs: ", recommendedSongs)

	aiRecommendedSongs := AIRecommendSongs(song, 5)
	fmt.Println("Recommended songs with LLM: ", aiRecommendedSongs)

	songs := recommendedSongs.Songs
	randomSong := songs[0]

	AddSongToQueue(randomSong)
	fmt.Println("Song added to queue:", randomSong)

	time.Sleep(5 * time.Second)

	SkipCurrentSong()
}
