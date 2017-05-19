package musicserver

import (
	"html/template"
	"net/http"
	"strings"
	"os"
	"io"

	"./playlist"
)

func homeHandler(w http.ResponseWriter, req *http.Request) {
	// Check if alias set for this ip
	if _, aliasSet := al.Alias(ip(req.RemoteAddr)); !aliasSet {
		http.Redirect(w, req, url("/alias"), http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("./templates/home.html")
	tmpl.Execute(w, newPageInfo(req.RemoteAddr))
}

func aliasHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tmpl, _ := template.ParseFiles("./templates/alias.html")
		tmpl.Execute(w, nil)
		return
	}

	newAlias := strings.TrimSpace(req.PostFormValue("alias_value"))
	if len(newAlias) < 1 {
		http.Redirect(w, req, url("/alias"), http.StatusSeeOther)
		return
	}

	// Set alias in the manager
	al.SetAlias(ip(req.RemoteAddr), newAlias)
	// Update listed aliases in the playlist in new goroutine
	go pl.UpdateAlias(ip(req.RemoteAddr), newAlias)

	http.Redirect(w, req, url("/"), http.StatusSeeOther)
}

func queueVideoHandler(w http.ResponseWriter, req *http.Request) {
	ip := ip(req.RemoteAddr)
	if req.Method != http.MethodPost {
		http.Redirect(w, req, url("/"), http.StatusSeeOther)
		return
	}
	// Check if alias set for this ip
	alias, aliasSet := al.Alias(ip)
	if !aliasSet {
		http.Redirect(w, req, url("/alias"), http.StatusSeeOther)
		return
	}

	if !pl.Available(ip) {
		tmpl, _ := template.ParseFiles("./templates/not_added.html")
		tmpl.Execute(w, "You have the maxium amount of videos queued.")
		return
	}

	newLink := req.PostFormValue("video_link")
	newVideo := playlist.NewVideo(ip, alias)
	pl.AddVideo(newVideo)
	go downloadVideo(newLink, newVideo.UUID)

	tmpl, _ := template.ParseFiles("./templates/added.html")
	tmpl.Execute(w, nil)
}

func uploadVideoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, url("/"), http.StatusSeeOther)
		return
	}

	// Open video file from post request before redirects can occur or the 
	// connection may be reset.
	file, header, err := req.FormFile("video_file")
	defer file.Close()
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/not_added.html")
		tmpl.Execute(w, "Cannot parse uploaded file")
		return
	}

	ip := ip(req.RemoteAddr)
	// Check if alias set for this ip
	alias, aliasSet := al.Alias(ip)
	if !aliasSet {
		http.Redirect(w, req, url("/alias"), http.StatusSeeOther)
		return
	}

	if !pl.Available(ip) {
		tmpl, _ := template.ParseFiles("./templates/not_added.html")
		tmpl.Execute(w, "You have the maximum amount of videos queued")
		return
	}

	newVid := playlist.NewVideo(ip, alias)
	// Gen file path with filename as uuid and get file extension from header
	newPath := vidFolder + "/" + newVid.UUID  + fileExt(header.Filename)

	// Create new file
	newFile, err := os.Create(newPath)
	defer newFile.Close()
	if err != nil {
		tmpl, _ := template.ParseFiles("./templates/not_added.html")
		tmpl.Execute(w, err.Error())
		return
	}

	// Write file
	_, err = io.Copy(newFile, file)
	if err != nil {
		tmpl, _ := template.ParseFiles("./templates/not_added.html")
		tmpl.Execute(w, err.Error())
		return
	}

	// Add information to Video struct
	newVid.Title = stripFileExt(header.Filename)
	newVid.File = newPath
	newVid.Ready = true

	// Add to playlist
	pl.AddVideo(newVid)

	tmpl, _ := template.ParseFiles("templates/added.html")
	tmpl.Execute(w, nil)
}

func userRemoveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, url("/"), http.StatusSeeOther)
		return
	}

	uuid := req.PostFormValue("video_id")
	if pl.VideoIP(uuid) == ip(req.RemoteAddr) {
		pl.RemoveVideo(uuid)
	}
	http.Redirect(w, req, url("/"), http.StatusSeeOther)
}
