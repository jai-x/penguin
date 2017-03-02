package musicserver

import (
	"../help"
	"strings"
	"net/http"
	"html/template"
	"encoding/json"
)

// Return homepage
func homeHandler(w http.ResponseWriter, req *http.Request) {
	_, aliasExists := Q.GetAlias(req.RemoteAddr)

	if !aliasExists {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-type", "text/html")

		plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
		homeTemplate, _ := template.ParseFiles("templates/home.html")
		homeTemplate.Execute(w, plInfo)
	}
}

// Endpoint for setting alias and returns alias set page
func aliasHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		newAlias := req.PostFormValue("alias_value")

		// Check of alias is whitespace
		if len(strings.TrimSpace(newAlias)) == 0 {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
		}

		Q.SetAlias(req.RemoteAddr, newAlias)

		http.Redirect(w, req, "/", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-type", "text/html")

		aliasTemplate, _ := template.ParseFiles("templates/alias.html")
		aliasTemplate.Execute(w, nil)
	}
}

// Endpoint for queuing videos via link
func queueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		_, aliasExists := Q.GetAlias(req.RemoteAddr)

		if !aliasExists {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
			return
		}

		videoLink := req.PostFormValue("video_link")

		// Submitted video link is blank
		if len(strings.TrimSpace(videoLink)) == 0 {
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		// If user has max added videos
		if !Q.CanAddVideo(req.RemoteAddr) {
			vidNotAddedTempl, _ := template.ParseFiles("templates/not_added.html")
			vidNotAddedTempl.Execute(w, nil)
			return
		}

		// Started in new go routine to prevent response waiting
		go Q.DownloadAndAddVideo(req.RemoteAddr, videoLink)

		vidAddedTempl, _ := template.ParseFiles("templates/added.html")
		vidAddedTempl.Execute(w, nil)
	} else {
		// Redirect back to homepage if not a POST request)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

func userRemoveHandler(w http.ResponseWriter, req *http.Request) {
	// Get video id from post data
	if req.Method == http.MethodPost {
		id := req.PostFormValue("video_id")
		ip := help.GetIP(req.RemoteAddr)
		Q.UserRemoveVideo(id, ip)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Endpoint to return playlist JSON
func playlistHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	info := Q.GetPlaylistInfo(req.RemoteAddr)
	json.NewEncoder(w).Encode(info)
}