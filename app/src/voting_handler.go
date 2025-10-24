package main

import (
	"sync"
	"math/rand"
)


type SongWithVotes struct {
	song Song
	votes map[string]bool
}

type VoteTally struct {
	songs     map[string]*SongWithVotes
	userVotes map[string]string
	mu        sync.Mutex
}


func initVotes(recommendedSongsList RecommendedSongs) *VoteTally {

	songsToVoteFor := &VoteTally{
		songs: make(map[string]*SongWithVotes),
		userVotes: make(map[string]string),
	}

	for _, song := range recommendedSongsList.songs {
		songsToVoteFor.songs[song.trackId] = &SongWithVotes{
			song: song,
			votes: make(map[string]bool),
		}
	}

	return songsToVoteFor
}


func (songsToVoteFor *VoteTally) Add(songID, uuid string) {

	songsToVoteFor.mu.Lock()
	defer songsToVoteFor.mu.Unlock()

	if oldSongID, ok := songsToVoteFor.userVotes[uuid]; ok {
		delete(songsToVoteFor.songs[oldSongID].votes, uuid)
	}

	if song, ok := songsToVoteFor.songs[songID]; ok {
		song.votes[uuid] = true
		songsToVoteFor.userVotes[uuid] = songID
	}
}


func (songsToVoteFor *VoteTally) Count(songID string) int {

	if song, ok := songsToVoteFor.songs[songID]; ok {
		return len(song.votes)
	}
	return 0
}


func (songsToVoteFor *VoteTally) DetermineWinner() Song {

	songsToVoteFor.mu.Lock()
	defer songsToVoteFor.mu.Unlock()

	var topSongs []*Song
	maxVotes := -1

	for _, songWithVotes := range songsToVoteFor.songs {
		voteCount := len(songWithVotes.votes)

		if voteCount > maxVotes {
			maxVotes = voteCount
			topSongs = []*Song{&songWithVotes.song}
		} else if voteCount == maxVotes {
			topSongs = append(topSongs, &songWithVotes.song)
		}
	}

	maxIndex := len(topSongs) - 1

	if maxIndex == 0 {
		topSong := topSongs[0]
		return *topSong
	} else {
		randomIndex := rand.Intn(maxIndex)
		randomTopSong := topSongs[randomIndex]
		return *randomTopSong
	}

}
