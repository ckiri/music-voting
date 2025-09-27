package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type VoteMessage struct {
	Vote string `json:"vote"`
}

type StateMessage struct {
	Counts map[string]int `json:"counts"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Voting state
	counts = map[string]int{
		"box1": 0,
		"box2": 0,
		"box3": 0,
	}
	mu sync.Mutex

	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Parse the template fresh on every request (so live reload works)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if err := tmpl.Execute(w, counts); err != nil {
		log.Println("Template execute error:", err)
	}
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Register client
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	// Send initial state
	sendState()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var vote VoteMessage
		if err := json.Unmarshal(msg, &vote); err == nil {
			mu.Lock()
			counts[vote.Vote]++
			mu.Unlock()
			sendState()
		}
	}

	// Remove client
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
}

func sendState() {
	mu.Lock()
	state := StateMessage{Counts: counts}
	mu.Unlock()

	data, _ := json.Marshal(state)

	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWS)

	// Serve static files (CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
