package server

// Sets urls and defines handler functions for each url

import (
	"log"
	"net/http"
	"strconv"

	"../config"
)

var (
	port string
	vidFolder string
)

func Init() {
	log.Println("Server init...")
	// Port is required as a string with a prefix colon
	port = ":" + strconv.Itoa(config.Config.Port)
	vidFolder = config.Config.DownloadFolder
}

func Run() {
	// Url Handlers
	// When a url is called, it spawns a new goroutine that runs the specified handler function

	// Debug url endpoints
	http.HandleFunc("/debug/list", debugListHandler)
	http.HandleFunc("/debug/np", debugNPHandler)
	// AJAX handlers
	http.HandleFunc("/ajax/queue", ajaxQueueHandler)
	http.HandleFunc("/ajax/upload", ajaxUploadHandler)
	// Admin url endpoints
	http.HandleFunc("/admin/logout", adminLogoutHandler)
	http.HandleFunc("/admin/login", adminLoginHandler)
	http.HandleFunc("/admin/remove", adminRemoveHandler)
	http.HandleFunc("/admin", adminHandler)
	// Regular url endpoints
	http.HandleFunc("/remove", removeHandler)
	http.HandleFunc("/alias", aliasHandler)
	http.HandleFunc("/queue", queueHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", homeHandler)
	// Static webpage files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// Video files
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(vidFolder))))

	// Run the server
	log.Println("Running music server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
