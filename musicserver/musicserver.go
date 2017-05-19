package musicserver

import (
	"log"
	"net/http"

	"./admin"
	"./alias"
	"./playlist"
	"./youtube"
)

var (
	al alias.AliasMgr
	ad admin.AdminSessions
	pl playlist.Playlist

	vidFolder string = "/tmp/penguin"
)

func Run() {
	// Create new instances of the main strucs
	al = alias.NewAliasMgr()
	ad = admin.NewAdminSessions("password", false)
	pl = playlist.NewPlaylist(3)

	// Set youtube settings
	youtube.SetYTDLPath("./dist/youtube-dl")
	youtube.SetFFMPEGPath("./dist/ffmpeg")
	youtube.SetDownloadFolder(vidFolder)

	// User facing url endpoints
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/alias", aliasHandler)
	http.HandleFunc("/queue", queueVideoHandler)
	http.HandleFunc("/upload", uploadVideoHandler)
	http.HandleFunc("/remove", userRemoveHandler)
	// Admin url endpoints
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/admin/login", adminLoginHandler)
	http.HandleFunc("/admin/logout", adminLogoutHandler)
	http.HandleFunc("/admin/remove", adminRemoveHandler)
	// Serve website static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Serve downloaded media files
	ms := http.FileServer(http.Dir(vidFolder))
	http.Handle("/media/", http.StripPrefix("/media/", ms))

	// Start server
	log.Println("Serving on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
