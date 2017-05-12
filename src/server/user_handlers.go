package server

import (
	"os"
	"io"
	"strings"
	"net/http"
	"html/template"

	"../help"
	"../playlist"
)

// Return homepage
func homeHandler(w http.ResponseWriter, req *http.Request) {
	_, aliasExists := playlist.GetAlias(req.RemoteAddr)

	if !aliasExists {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
	} else {
		tmpl, _ := template.ParseFiles("templates/home.html")
		tmpl.Execute(w, fetchInfo(req.RemoteAddr))
	}
}

// Endpoint for setting alias and returns alias set page
func aliasHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tmpl, _ := template.ParseFiles("templates/alias.html")
		tmpl.Execute(w, nil)
		return
	}

	newAlias := req.PostFormValue("alias_value")

	// Check if alias is whitespace
	if len(strings.TrimSpace(newAlias)) == 0 {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
	}
	playlist.SetAlias(req.RemoteAddr, newAlias)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func queueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	_, aliasExists := playlist.GetAlias(req.RemoteAddr)

	// Trying to upload with no alias set
	if !aliasExists {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
		return
	}

	if !playlist.CanAddVideo(req.RemoteAddr) {
		tmpl, _ := template.ParseFiles("templates/not_added.html")
		tmpl.Execute(w, "You already have the max number of videos queued.")
		return
	}

	// Get video link from form post value and add to playlist
	vidLink := req.PostFormValue("video_link")
	playlist.AddVideoLink(req.RemoteAddr, vidLink)

	// Show link added page
	tmpl, _ := template.ParseFiles("templates/added.html")
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// Must open file from post before else or connection is reset
	file, header, err := req.FormFile("video_file")
	defer file.Close()
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/not_added.html")
		tmpl.Execute(w, "Cannot parse uploaded file")
		return
	}

	// Check for set user alias for ip
	alias, aliasExists := playlist.GetAlias(req.RemoteAddr)
	if !aliasExists {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
		return
	}

	// Check for number of videos added
	if !playlist.CanAddVideo(req.RemoteAddr) {
		tmpl, _ := template.ParseFiles("templates/not_added.html")
		tmpl.Execute(w, "You already have the max number of videos queued.")
		return
	}

	// Gen uuid and path for new file
	uuid := help.GenUUID()
	path := vidFolder + "/" + uuid

	// Create file
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/not_added.html")
		tmpl.Execute(w, err.Error)
		return
	}

	// Write file
	_, err = io.Copy(out, file)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/not_added.html")
		tmpl.Execute(w, err.Error)
		return
	}

	// Struct for new video
	newVid := playlist.Video{
		UUID: uuid,
		Title: header.Filename,
		File: path,
		IpAddr: help.GetIP(req.RemoteAddr),
		Alias: alias,
		Ready: true,
		Played: false,
	}

	go playlist.AddVideoStruct(newVid)

	tmpl, _ := template.ParseFiles("templates/added.html")
	tmpl.Execute(w, nil)
}

func removeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		remUUID := req.PostFormValue("video_id")

		if playlist.AddrOwnsVideo(req.RemoteAddr, remUUID) {
			playlist.RemoveVideo(remUUID)
		}
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
