package musicserver

// Main package for music server
// Sets urls and defines handler functions for each url

import (
	"log"
	"net/http"

	"../state"
	"../admin"
	"../config"
)

var (
	Q state.ProcessQueue
	A admin.AdminSessions
)

func Init(configPath string) {
	config.Init(configPath)

	Q.Init()
	A.Init()

	config.End()
}

func Run() {
	// Url Handlers
	// When a url is called, it spawns a new goroutine that runs the specifeed handler function

	// Admin url endpoints
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/admin/login", adminLoginHandler)
	http.HandleFunc("/admin/logout", adminLogoutHandler)
	http.HandleFunc("/admin/kill", adminKillHandler)
	http.HandleFunc("/admin/remove", adminRemoveHandler)
	// Regular url endpoints
	http.HandleFunc("/alias", aliasHandler)
	http.HandleFunc("/queue", queueHandler)
	http.HandleFunc("/remove", userRemoveHandler)
	http.HandleFunc("/", homeHandler)
	// AJAX Endpoints
	http.HandleFunc("/ajax/queue", ajaxQueueHandler)
	http.HandleFunc("/ajax/playlist", ajaxPlaylistHandler)
	http.HandleFunc("/ajax/adminplaylist", ajaxAdminPlaylistHandler)
	// Static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start video player service in a separate goroutine
	go Q.VideoPlayerService()


	// Run the server
	log.Println("Running music server on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
