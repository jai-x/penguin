package musicserver

import (
	"flag"
	"log"
	"net/http"
	"time"

	"./admin"
	"./alias"
	"./config"
	"./player"
	"./playlist"
	"./templatecache"
)

var (
	al *alias.AliasMgr
	ad *admin.AdminSessions
	pl *playlist.Playlist
	vd *player.VideoPlayer
	tl *templatecache.TmplCache
	conf *config.Config

	// Command line flag set variables
	playVideos   bool
	useTmplCache bool
)

func Init() {
	// Parse command line
	flag.BoolVar(&playVideos, "play-videos", true, "Set video playback behaviour")
	flag.BoolVar(&useTmplCache, "template-cache", true, "Set template caching behaviour")
	flag.Parse()

	var err error
	conf, err = config.ReadConfig()
	if err != nil {
		log.Fatalln("Cannot read config file:", err.Error())
	}
	// Create new instances of the main strucs
	al = alias.NewAliasMgr()
	ad = admin.NewAdminSessions(conf.AdminPass, false)
	pl = playlist.NewPlaylist(conf.Buckets)

	// Set domain so that templates have the correct absolute hyperlinks
	templatecache.SetDomain(conf.ServerDomain)
	tl = templatecache.NewTemplateCache(conf.TemplateDir, useTmplCache)
	if !useTmplCache {
		log.Println("Template caching disabled")
	}
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
	ms := http.FileServer(http.Dir(conf.VidFolder))
	http.Handle("/media/", http.StripPrefix("/media/", ms))

	// Start video player
	if playVideos {
		go videoPlayer()
	} else {
		log.Println("Video player disabled")
	}

	// Start server
	log.Println("Serving on port:", conf.Port)
	err := http.ListenAndServe(conf.Port, nil)
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
			vd = player.NewVideoPlayer(conf.VidTimout, conf.VidExe, conf.VidArgs, newVid)
			vd.Play()
		}
	}
}
