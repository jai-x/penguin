package musicserver

import (
	"log"
	"net/http"
	"time"
	"flag"

	"./admin"
	"./alias"
	"./playlist"
	"./youtube"
	"./player"
	"./templatecache"
)

var (
	al alias.AliasMgr
	ad admin.AdminSessions
	pl playlist.Playlist
	vd player.VideoPlayer
	tl templatecache.TmplCache


	vidFolder    string = "/tmp/penguin"
	vidPlayer    string = "mpv"
	vidArgs      string = "-fs"
	vidTimout    string = "547s"
	adminPass    string = "password"
	serverDomain string = ""
	templateDir  string = "./templates"
	ytdlExe      string = "./dist/youtube-dl"
	ffmpegExe    string = "./dist/ffmpeg"
	plBuckets    int = 5
)

func Run() {
	// Create new instances of the main strucs
	al = alias.NewAliasMgr()
	ad = admin.NewAdminSessions(adminPass, false)
	pl = playlist.NewPlaylist(plBuckets)
	tl = templatecache.NewTemplateCache(templateDir, serverDomain)

	// Set youtube settings
	youtube.SetYTDLPath(ytdlExe)
	// Optional set
	youtube.SetFFMPEGPath(ffmpegExe)
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
	http.HandleFunc("/admin/kill", adminKillVideoHandler)

	// AJAX url endpoints
	http.HandleFunc("/ajax/queue", ajaxQueueHandler)
	http.HandleFunc("/ajax/upload", ajaxUploadHandler)
	http.HandleFunc("/ajax/playlist", ajaxPlaylistHandler)
	http.HandleFunc("/ajax/admin/playlist", ajaxAdminPlaylistHandler)

	// Debug url endpoints
	http.HandleFunc("/debug/playlist", debugListHandler)
	http.HandleFunc("/debug/ip", debugIPHandler)

	// Serve website static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve downloaded media files
	ms := http.FileServer(http.Dir(vidFolder))
	http.Handle("/media/", http.StripPrefix("/media/", ms))

	var noPlayer bool
	flag.BoolVar(&noPlayer, "noplayer", false, "Disables video playback")
	flag.Parse()

	// Start video player
	if !noPlayer {
		go videoPlayer()
	}

	// Start server
	log.Println("Serving on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func videoPlayer() {
	emptyVid := playlist.Video{}
	for {
		newVid := pl.NextVideo()
		if newVid == emptyVid {
			log.Println(`(/'-')/ No Videos to Play \('-'\)`)
			time.Sleep(2 * time.Second)
		} else {
			vd = player.NewVideoPlayer(vidPlayer, vidArgs, newVid.File, vidTimout)
			vd.Play()
		}
	}
}
