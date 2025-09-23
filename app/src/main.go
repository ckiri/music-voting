package main

import (
	//"fmt"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	InitSpotify()
	playbackState := GetPlaybackState()
	song := playbackState.song
	fmt.Println("Song playing right now:", song)
	recommendedSongs := AIRecommendSongs(song)
	fmt.Println(recommendedSongs)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
