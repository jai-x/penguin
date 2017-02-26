package mserver

import (
	"strings"
	"net/http"
	"html/template"
)

func homeHandler(w http.ResponseWriter, req *http.Request) {
	userAlias, aliasExists := Q.getAliasFromAddress(req.RemoteAddr)

	if !aliasExists {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
	} else {
		plInfo := getPlaylistInfo()
		pageInfo := PageInfo{plInfo.Playlist, plInfo.NowPlaying, userAlias}
		homeTemplate, _ := template.ParseFiles("templates/home.html")
		homeTemplate.Execute(w, pageInfo)
	}
}

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
		go Q.setNewAlias(req.RemoteAddr, newAlias)

		http.Redirect(w, req, "/home", http.StatusSeeOther)
	} else {
		aliasTemplate, _ := template.ParseFiles("templates/alias.html")
		aliasTemplate.Execute(w, nil)
	}
}

func queueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {

		_, aliasExists := Q.getAliasFromAddress(req.RemoteAddr)

		if !aliasExists {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
		} else {
			req.ParseForm()
			videoLink := req.Form["video_link"][0]

			go Q.downloadAndAddVideo(req.RemoteAddr, videoLink)

			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	} else {
		// Redirect back to homepage if not a POST request)
		http.Redirect(w, req, "/home", http.StatusSeeOther)
	}
}