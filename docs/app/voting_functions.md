# Available functions to handle voting for songs

## Init

Once the steps from [spotify_api](./spotify_api.md) have been followed and a `RecommendedSongs {}` struct has been recived the following functions and structs are given to handle the logic behind voting for songs.

A Wrapper-like object which contains a song and a list of users who have voted for the song:
```go
type SongWithVotes struct {
	song   Song
	votes  map[string]bool
}
```

A stuct witch contains the information of all songs and given users votes:
```go
type VoteTally struct {
	songs     map[string]*SongWithVotes
	userVotes map[string]string
	mu        sync.Mutex
}
```

A init function which allows the user to start a voting session. The returned `VoteTally` object contains all neccecary functions for voting, reading and selecting the most voted for song:
```go
func initVotes(recommendedSongsList RecommendedSongs) *VoteTally { ... }
```

## Using the voting session

Once the session has been initialised, the following functions are given for use:

`.Add(songId, userId)` can be called to add a users vote for a song. If the user already has voted for a different song, the old vote is redacted:
```go
func (songsToVoteFor *VoteTally) Add(songID, uuid string) { ... }
```

`.Count(songId)` can be called to return the total count of votes for a single song:
```go
func (songsToVoteFor *VoteTally) Count(songID string) int { ... }
```

`.DetermineWinner()` can be called to return a song object which has the most amount of votes:
```go
func (songsToVoteFor *VoteTally) DetermineWinner() Song { ... }
```

## Bugs/ToDo:
- Make decision on how users are identified in the programm and apply this to userId related functions
- Handle empty recommendation list (?)
