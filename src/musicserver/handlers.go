package musicserver

import (
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
		plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
		homeTemplate, _ := template.ParseFiles("templates/home.html")
		homeTemplate.Execute(w, plInfo)
	}
}

// Endpoint for setting alias and returns alias set page
func aliasHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		newAlias := req.Form["alias_value"][0]

		if len(strings.TrimSpace(newAlias)) == 0 {
			// Alias is just whitespace
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
			return
		}

		// Using a new goroutine prevents waiting on mutex for http response
		go Q.SetAlias(req.RemoteAddr, newAlias)

		http.Redirect(w, req, "/home", http.StatusSeeOther)
	} else {
		aliasTemplate, _ := template.ParseFiles("templates/alias.html")
		aliasTemplate.Execute(w, nil)
	}
}

// Endpoint for queuing videos via link
func queueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {

		_, aliasExists := Q.GetAlias(req.RemoteAddr)

		if !aliasExists {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
		} else {
			req.ParseForm()
			videoLink := req.Form["video_link"][0]

			go Q.DownloadAndAddVideo(req.RemoteAddr, videoLink)

			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	} else {
		// Redirect back to homepage if not a POST request)
		http.Redirect(w, req, "/home", http.StatusSeeOther)
	}
}

// Endpoint to return playlist JSON
func playlistHandler(w http.ResponseWriter, req *http.Request) {
	info := Q.GetPlaylistInfo(req.RemoteAddr)
	json.NewEncoder(w).Encode(info)
}