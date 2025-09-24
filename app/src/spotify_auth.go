package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Defining a workflow for the auth / init process
// Function InitSpotify() is called on runtime
// Following checks are made:
// 1. Is there already a token.json file? Yes --> 1.1 No --> 2
// 1.1 Is the token expired? Yes --> Refresh No --> Return the token for further auth
// 2. No token is present and the auth workflow is started, user needs to copy link in stdout to auth your key

var (
	spotifyAuth *oauth2.Config
)

const (
	TokenFile = "token.json"
)

func loadEnv() {
	/*
	 * This function loads the .env file. TODO: Find a better way to reference the file
	 */
	err := godotenv.Load("../env/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func initSpotifyAuth() {
	/*
	 * Initialises the object used for creating a token
	 */
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
	/*
	 * Returns the auth url needed for authenticating the user
	 */
	return spotifyAuth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func exchangeToken(code string) (*oauth2.Token, error) {
	token, err := spotifyAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func refreshTokenWithRefreshToken(refreshToken string) (*oauth2.Token, error) {
	tokenSource := spotifyAuth.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func saveToken(path string, token *oauth2.Token) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}

func ReadToken(path string) (*oauth2.Token, error) {
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
	/*
	 * This func ensures that the token.json file always has a valid token defined
	 */

	// Initialise the structures
	loadEnv()
	initSpotifyAuth()

	// Try to read the token, check if valid and refresh if possible
	token, err := ReadToken(TokenFile)
	if err == nil && token.Valid() {
		fmt.Println("Token found and is valid. No need to re-authorize.")
		return
	} else if err == nil && !token.Valid() {
		fmt.Printf("Expired Token was found, refreshing...")
		refreshToken := token.RefreshToken
		newToken, err := refreshTokenWithRefreshToken(refreshToken)
		if err == nil {
			saveToken(TokenFile, newToken)
			token = newToken
			return
		}
		fmt.Printf("Error refreshing token: %v\n", err)
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

		saveToken(TokenFile, token)
		tokenChan <- token

		fmt.Fprintf(w, "Access Token received. You can close this window.")
	})

	fmt.Println("Open the following URL in your browser to authorize:")
	fmt.Println(getAuthURL())

	// Start the server in a separate goroutine so it doesn't block
	// Here the value from redirect uri in .env is taken split, part after the : port is taken and split
	// again to extract only the port number. This is then given to the served http
	go func() {
		port := strings.Split(strings.Split(os.Getenv("SPOTIFY_REDIRECT_URI"), ":")[2], "/")[0]
		print(port)
		log.Fatal(http.ListenAndServe(":"+ port, nil))
	}()

	token = <-tokenChan
}
