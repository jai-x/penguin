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

	// Config file variables
	vidFolder    string = "/tmp/penguin"
	vidExe       string = "mpv"
	vidArgs      string = "-fs"
	vidTimout    string = "547s"
	adminPass    string = "password"
	serverDomain string = "http://192.168.0.15:8080"
	templateDir  string = "./templates"
	ytdlExe      string = "./dist/youtube-dl"
	ffmpegExe    string = "./dist/ffmpeg"
	port         string = ":8080"
	plBuckets    int = 5

	// Command line flag set variables
	playVideos   bool
	useTmplCache bool
)

func Init() {
	// Parse command line
	flag.BoolVar(&playVideos, "play-videos", true, "Set video playback behaviour")
	flag.BoolVar(&useTmplCache, "template-cache", true, "Set template caching behaviour")
	flag.Parse()

	// Create new instances of the main strucs
	al = alias.NewAliasMgr()
	ad = admin.NewAdminSessions(adminPass, false)
	pl = playlist.NewPlaylist(plBuckets)

	// Set domain so that templates have the correct absolute hyperlinks
	templatecache.SetDomain(serverDomain)
	tl = templatecache.NewTemplateCache(templateDir, useTmplCache)
	if !useTmplCache {
		log.Println("Template caching disabled")
	}

	// Set video player setings
	player.SetPlayer(vidExe, vidArgs)
	player.SetTimeout(vidTimout)

	// Set youtube settings
	youtube.SetYTDLPath(ytdlExe)
	// Optional set
	youtube.SetFFMPEGPath(ffmpegExe)
	youtube.SetDownloadFolder(vidFolder)
}

func Run() {
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
	http.HandleFunc("/debug/header", debugHeaderHandler)

	// Serve website static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve downloaded media files
	ms := http.FileServer(http.Dir(vidFolder))
	http.Handle("/media/", http.StripPrefix("/media/", ms))

	// Start video player
	if playVideos {
		go videoPlayer()
	} else {
		log.Println("Video player disabled")
	}

	// Start server
	log.Println("Serving on port:", port)
	err := http.ListenAndServe(port, nil)
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
			vd = player.NewVideoPlayer(newVid)
			vd.Play()
		}
	}
}
