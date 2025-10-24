package main

import (
	"testing"
)

func makeSongs() RecommendedSongs {
	return RecommendedSongs{
		songs: []Song{
			{trackId: "1", trackName: "Song A", artistName: "Artist A"},
			{trackId: "2", trackName: "Song B", artistName: "Artist B"},
			{trackId: "3", trackName: "Song C", artistName: "Artist C"},
		},
	}
}

func TestInitVotes(t *testing.T) {
	recommendedSongs := makeSongs()
	voteTally := initVotes(recommendedSongs)

	if len(voteTally.songs) != 3 {
		t.Errorf("expected 3 songs, got %d", len(voteTally.songs))
	}

	for _, s := range recommendedSongs.songs {
		if _, ok := voteTally.songs[s.trackId]; !ok {
			t.Errorf("song %s not initialized in VoteTally", s.trackId)
		}
	}
}

func TestAddAndCount(t *testing.T) {
	recommendedSongs := makeSongs()
	voteTally := initVotes(recommendedSongs)

	voteTally.Add("1", "user1")
	voteTally.Add("1", "user2")
	voteTally.Add("2", "user3")

	if got := voteTally.Count("1"); got != 2 {
		t.Errorf("expected 2 votes for song 1, got %d", got)
	}
	if got := voteTally.Count("2"); got != 1 {
		t.Errorf("expected 1 vote for song 2, got %d", got)
	}
	if got := voteTally.Count("3"); got != 0 {
		t.Errorf("expected 0 votes for song 3, got %d", got)
	}
}

func TestChangeVote(t *testing.T) {
	recommendedSongs := makeSongs()
	voteTally := initVotes(recommendedSongs)

	voteTally.Add("1", "user1")
	voteTally.Add("2", "user1") // change vote

	if got := voteTally.Count("1"); got != 0 {
		t.Errorf("expected 0 votes for song 1 after change, got %d", got)
	}
	if got := voteTally.Count("2"); got != 1 {
		t.Errorf("expected 1 vote for song 2 after change, got %d", got)
	}
}

func TestDetermineWinnerClearWinner(t *testing.T) {
	recommendedSongs := makeSongs()
	voteTally := initVotes(recommendedSongs)

	voteTally.Add("1", "user1")
	voteTally.Add("1", "user2")
	voteTally.Add("2", "user3")

	winner := voteTally.DetermineWinner()

	if winner.trackId != "1" {
		t.Errorf("expected song 1 as winner, got %s", winner.trackId)
	}
}

func TestDetermineWinnerTie(t *testing.T) {

	recommendedSongs := makeSongs()
	voteTally := initVotes(recommendedSongs)

	voteTally.Add("1", "user1")
	voteTally.Add("2", "user2")

	winner := voteTally.DetermineWinner()
	if winner.trackId != "1" && winner.trackId != "2" {
		t.Errorf("expected winner to be 1 or 2, got %s", winner.trackId)
	}
}
