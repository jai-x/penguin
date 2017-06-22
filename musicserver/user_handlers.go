package musicserver

import (
	"net/http"
	"strings"
)

func homeHandler(w http.ResponseWriter, req *http.Request) {
	// Check if alias set for this ip
	ip := getIPFromRequest(req)
	if _, aliasSet := al.Alias(ip); !aliasSet {
		http.Redirect(w, req, url("/alias"), http.StatusSeeOther)
		return
	}
	tl.Render(w, "home", newPlaylistInfo(ip))
}

func aliasHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tl.Render(w, "alias", nil)
		return
	}

	newAlias := strings.TrimSpace(req.PostFormValue("alias_value"))
	if len(newAlias) < 1 {
		http.Redirect(w, req, url("/alias"), http.StatusSeeOther)
		return
	}

	ip := getIPFromRequest(req)
	// Set alias in the manager
	al.SetAlias(ip, newAlias)
	// Update listed aliases in the playlist in new goroutine
	go pl.UpdateAlias(ip, newAlias)

	http.Redirect(w, req, url("/"), http.StatusSeeOther)
}

func queueVideoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, url("/"), http.StatusSeeOther)
		return
	}

	if err := queueLink(req); err != nil {
		tl.Render(w, "not_added", err.Error())
		return
	}

	tl.Render(w, "added", nil)
}

func uploadVideoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, url("/"), http.StatusSeeOther)
		return
	}

	if err := queueUploadedVideo(req); err != nil {
		tl.Render(w, "not_added", err.Error())
		return
	}

	tl.Render(w, "added", nil)
}

func userRemoveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, url("/"), http.StatusSeeOther)
		return
	}

	ip := getIPFromRequest(req)
	uuid := req.PostFormValue("video_id")
	if pl.VideoIP(uuid) == ip {
		pl.RemoveVideo(uuid)
	}
	http.Redirect(w, req, url("/"), http.StatusSeeOther)
}
