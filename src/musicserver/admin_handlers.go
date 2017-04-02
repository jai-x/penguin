package musicserver

import (
	"log"
	"time"
	"net/http"
	"os/exec"

	"../templatecache"
)

func adminHandler(w http.ResponseWriter, req *http.Request) {
	if !A.ValidSession(req.RemoteAddr) {
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
	}

	plInfo := Q.GetPlaylistInfo(req.RemoteAddr)

	templatecache.Render(w, "admin", &plInfo)
}

func adminLoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		correct := A.VerifyPassword(req.PostFormValue("admin_pwd"))

		if !correct {
			templatecache.Render(w, "login_incorrect", nil)
			return
		} else {
			A.StartSession(req.RemoteAddr)
			http.Redirect(w, req, "/admin", http.StatusSeeOther)
		}

	} else {
		templatecache.Render(w, "admin_login", nil)
	}
}

func adminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	A.EndSession(req.RemoteAddr)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func adminKillHandler(w http.ResponseWriter, req *http.Request) {
	if !A.ValidSession(req.RemoteAddr) {
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}
	// Use killall to kill music players
	killPlayer := exec.Command("killall", "mpv")
	killPlayer.Run()

	log.Println("Admin killed current video")

	// Allow for the video playerservice to cycle to next video
	// Only so that when the page refreshes it shows the video not playing
	time.Sleep(500 * time.Millisecond)

	// Go back to admin page
	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}

func adminRemoveHandler(w http.ResponseWriter, req *http.Request) {
	if !A.ValidSession(req.RemoteAddr) {
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}

	// Get video id from post data
	if req.Method == http.MethodPost {
		id := req.PostFormValue("video_id")
		Q.AdminRemoveVideo(id)
	}

	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}
