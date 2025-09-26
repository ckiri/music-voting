package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/gin-gonic/gin"
)

func main() {
	demoVoteTallying()
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}

func demoAddRandomRecommendedSongToQueueAndPlay() {

	recommendAmount := 5
	maxIndex := recommendAmount - 1

	InitSpotify()

	playbackState, err := GetPlaybackState()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Player state: ", playbackState)

	song := playbackState.song

	recommendedSongs, err := RecommendSongs(song, recommendAmount)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Recommended songs: ", recommendedSongs)

	aiRecommendedSongs, err := AIRecommendSongs(song, recommendAmount)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Recommended songs with LLM: ", aiRecommendedSongs)
	}

	songs := recommendedSongs.songs
	random := rand.Intn(maxIndex)
	randomSong := songs[random]

	err = AddSongToQueue(randomSong)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Song added to queue:", randomSong)
	}

	time.Sleep(5 * time.Second)

	SkipCurrentSong()
}

func demoVoteTallying() {

	recommendAmount := 5

	InitSpotify()

	playbackState, err := GetPlaybackState()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Player state: ", playbackState)

	song := playbackState.song

	recommendedSongs, err := RecommendSongs(song, recommendAmount)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Recommended songs: ", recommendedSongs)

	votingSession :=  initVotes(recommendedSongs)
	votingSession.Add(recommendedSongs.songs[0].trackId, "user-1")
	votingSession.Add(recommendedSongs.songs[0].trackId, "user-2")
	votingSession.Add(recommendedSongs.songs[0].trackId, "user-3")

	for _, song := range recommendedSongs.songs {
		fmt.Printf("%s has %d votes\n", song.trackName, votingSession.Count(song.trackId))
	}

	topVotedSong := votingSession.DetermineWinner()

	err = AddSongToQueue(topVotedSong)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Song added to queue:", topVotedSong)
	}

	time.Sleep(5 * time.Second)

	SkipCurrentSong()
}
