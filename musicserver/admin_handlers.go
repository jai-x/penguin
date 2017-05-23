package musicserver

import (
	"net/http"
	"time"
)

func adminHandler(w http.ResponseWriter, req *http.Request) {
	if ad.ValidSession(ip(req.RemoteAddr)) {
		tl.Render(w, "admin", newPlaylistInfo(req.RemoteAddr))
	} else {
		http.Redirect(w, req, url("/admin/login"), http.StatusFound)
	}
}

func adminLoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tl.Render(w, "admin_login", nil)
		return
	}

	pwd := req.PostFormValue("admin_pwd")
	if ad.ValidPassword(pwd) {
		ad.StartSession(ip(req.RemoteAddr))
		http.Redirect(w, req, url("/admin"), http.StatusSeeOther)
		return
	} else {
		tl.Render(w, "admin_bad_login", nil)
		return
	}
}

func adminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	ad.EndSession(ip(req.RemoteAddr))
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func adminRemoveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && ad.ValidSession(ip(req.RemoteAddr)) {
		remUUID := req.PostFormValue("video_id")
		pl.RemoveVideo(remUUID)
	}
	http.Redirect(w, req, url("/admin"), http.StatusSeeOther)
}

func adminKillVideoHandler(w http.ResponseWriter, req *http.Request) {
	if !ad.ValidSession(ip(req.RemoteAddr)) {
		http.Redirect(w, req, url("/admin/login"), http.StatusFound)
		return
	}

	vd.End()
	// Wait for playlist to cycle
	time.Sleep(500 * time.Millisecond)

	http.Redirect(w, req, url("/admin"), http.StatusFound)
}
