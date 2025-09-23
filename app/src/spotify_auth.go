package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	spotifyAuth *oauth2.Config
)

const tokenFile = "token.json"

func loadEnv() {
	err := godotenv.Load("../env/.env")
	if err != nil {
		log.Fatalf("Error loading .enf file: %v", err)
	}
}

func initSpotifyAuth() {
	spotifyAuth = &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		Scopes:       []string{"user-modify-playback-state", "user-read-playback-state"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
}

func getAuthURL() string {
	return spotifyAuth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func exchangeToken(code string) (*oauth2.Token, error) {
	token, err := spotifyAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving token to %s\n", path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}

func readToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func InitSpotify() {
	loadEnv()
	initSpotifyAuth()

	token, err := readToken(tokenFile)
	if err == nil && token.Valid() {
		fmt.Println("Token found and is valid. No need to re-authorize.")
		return
	}

	tokenChan := make(chan *oauth2.Token)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Fprintf(w, "Error: missing code parameter")
			return
		}

		token, err := exchangeToken(code)
		if err != nil {
			fmt.Fprintf(w, "Error exchanging token: %v", err)
			return
		}

		saveToken(tokenFile, token)
		tokenChan <- token

		fmt.Fprintf(w, "Access Token received. You can close this window.")
	})

	fmt.Println("Open the following URL in your browser to authorize:")
	fmt.Println(getAuthURL())

	// Start the server in a separate goroutine so it doesn't block
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	token = <-tokenChan

	fmt.Println("\n-------------------------------------")
	fmt.Println("Access Token received:")
	fmt.Printf("Token: %s\n", token.AccessToken)
	fmt.Printf("Refresh Token: %s\n", token.RefreshToken)
	fmt.Printf("Expires In: %v\n", token.Expiry)
	fmt.Println("-------------------------------------")
}
