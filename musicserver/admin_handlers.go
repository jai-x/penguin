package musicserver

import (
	"net/http"
	"time"
)

func adminHandler(w http.ResponseWriter, req *http.Request) {
	ip := getIPFromRequest(req)

	if ad.ValidSession(ip) {
		tl.Render(w, "admin", newPlaylistInfo(ip))
	} else {
		http.Redirect(w, req, url("/admin/login"), http.StatusFound)
	}
}

func adminLoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tl.Render(w, "admin_login", nil)
		return
	}

	ip := getIPFromRequest(req)
	pwd := req.PostFormValue("admin_pwd")
	if ad.ValidPassword(pwd) {
		ad.StartSession(ip)
		http.Redirect(w, req, url("/admin"), http.StatusSeeOther)
		return
	} else {
		tl.Render(w, "admin_bad_login", nil)
		return
	}
}

func adminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	ip := getIPFromRequest(req)
	ad.EndSession(ip)
	http.Redirect(w, req, url("/"), http.StatusSeeOther)
}

func adminRemoveHandler(w http.ResponseWriter, req *http.Request) {
	ip := getIPFromRequest(req)
	if req.Method == http.MethodPost && ad.ValidSession(ip) {
		remUUID := req.PostFormValue("video_id")
		pl.RemoveVideo(remUUID)
	}
	http.Redirect(w, req, url("/admin"), http.StatusSeeOther)
}

func adminKillVideoHandler(w http.ResponseWriter, req *http.Request) {
	ip := getIPFromRequest(req)
	if !ad.ValidSession(ip) {
		http.Redirect(w, req, url("/admin/login"), http.StatusFound)
		return
	}

	vd.End()
	// Wait for playlist to cycle
	time.Sleep(500 * time.Millisecond)

	http.Redirect(w, req, url("/admin"), http.StatusFound)
}
