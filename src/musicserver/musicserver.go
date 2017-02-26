package musicserver

import (
	"log"
	"net/http"

	"../state"
	"../admin"
)

var (
	Q state.Queue

	A admin.AdminSessions

	debugMode bool
)

func Init(debug bool) {
	// Set global debug mode
	debugMode = debug

	// Set timeout and debug 
	// Also intialise video state
	Q.Init(547, debugMode)
	// Intialise map of admin sessions
	A.Init("pass")

	// Debug check
	if debugMode {
		log.Println("####################")
		log.Println("#### DEBUG MODE ####")
		log.Println("####################")
	}
}

func Run() {
	// Url Handlers
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/admin/login", adminLoginHandler)
	http.HandleFunc("/admin/logout", adminLogoutHandler)
	http.HandleFunc("/admin/kill", adminKillHandler)

	http.HandleFunc("/playlist", playlistHandler)
	http.HandleFunc("/alias", aliasHandler)
	http.HandleFunc("/queue", queueHandler)
	http.HandleFunc("/home/", homeHandler)

	// Dont play videos in debug
	if !debugMode {
		// Start video player function in a separate goroutine
		go Q.PlayVideos()
	}

	// Run the server
	log.Println("Running music server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}