package musicserver

import (
	"log"
	"net/http"

	"../state"
	"../admin"
)

var (
	Q state.ProcessQueue
	A admin.AdminSessions

	debug bool
)

func Init(debug bool) {
	debug = debug
	// Set timeout and debug 
	// Also intialise video state
	Q.Init(547, 3, debug)
	// Intialise admin sessions with admin password
	A.Init("pass")

	// Debug check
	if debug {
		log.Println("####################")
		log.Println("#### DEBUG MODE ####")
		log.Println("####################")
	}
}

func Run(debug bool) {
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

	if !debug {
		// Start video player service in a separate goroutine
		go Q.VideoPlayerService()
	} else {
		// Debug endpoints 
		http.HandleFunc("/playlist", playlistHandler)
		log.Println("DEBUG: Additional endpoints: /playlist")
		http.HandleFunc("/rawplaylist", rawPlaylistHandler)
		log.Println("DEBUG: Additional endpoint: /rawplaylist")

		log.Println("DEBUG: VideoPlayerService suspended")
	}

	// Run the server
	log.Println("Running music server on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
