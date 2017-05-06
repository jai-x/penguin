package musicserver

// Main package for music server
// Sets urls and defines handler functions for each url

import (
	"log"
	"net/http"
	"strconv"

	"../state"
	"../admin"
	"../templatecache"
	"../config"
	"../help"
)

var (
	Q state.ProcessQueue
	A admin.AdminSessions
	port string
)

func Init(configPath string) {
	// Get config
	config.Init(configPath)

	Q.Init()
	A.Init()
	templatecache.Init()

	// Port is required as a string with a prefix colon
	port = ":" + strconv.Itoa(config.Config.Port)

	// Empty config
	config.End()

	help.PrintMasthead()
}

func Run() {
	// Url Handlers
	// When a url is called, it spawns a new goroutine that runs the specified handler function

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
	http.HandleFunc("/upload", fileUploadHandler)
	http.HandleFunc("/list", showList)
	http.HandleFunc("/", homeHandler)
	// AJAX Endpoints
	http.HandleFunc("/ajax/queue", ajaxQueueHandler)
	http.HandleFunc("/ajax/playlist", ajaxPlaylistHandler)
	http.HandleFunc("/ajax/adminplaylist", ajaxAdminPlaylistHandler)
	http.HandleFunc("/ajax/upload", ajaxFileUploadHandler)
	// Static file server
	http.Handle("/static/",http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start video player service in a separate goroutine
	go Q.VideoPlayerService()

	// Run the server
	log.Println("Running music server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
