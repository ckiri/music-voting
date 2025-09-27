package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

// Parse your template once (assuming templates/index.html exists)
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Parse fresh on every request
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	// Provide any data to the template
	data := map[string]interface{}{
		"Title": "Voting App",
		"box1":  0,
		"box2":  0,
		"box3":  0,
	}

	tmpl.Execute(w, data)
}

func livereloadWS(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	mu.Lock()
	clients[conn] = true
	mu.Unlock()
}

func notifyClients() {
	mu.Lock()
	defer mu.Unlock()
	for c := range clients {
		c.WriteMessage(websocket.TextMessage, []byte("reload"))
	}
}

func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch your template + CSS dirs
	watcher.Add("templates")
	watcher.Add("static/css")
	watcher.Add("static/js")

	for {
		select {
		case event := <-watcher.Events:
			log.Println("File change detected:", event)
			notifyClients()
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}

func main() {
	// Your main HTML page
	http.HandleFunc("/", handleHome)

	// Livereload endpoint
	http.HandleFunc("/livereload", livereloadWS)

	// Static files (CSS/JS/images)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start watcher in background
	go watchFiles()

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
