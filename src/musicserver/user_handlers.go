package musicserver

import (
	"os"
	"io"
	"fmt"
	"strings"
	"net/http"
	"html/template"

	"../help"
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

// Endpoint for queuing videos via webform
func queueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		_, aliasExists := Q.GetAlias(req.RemoteAddr)

		if !aliasExists {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
			return
		}

		// Get link from form
		videoLink := req.PostFormValue("video_link")

		// Submitted video link is blank
		if len(strings.TrimSpace(videoLink)) == 0 {
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		// If user has max added videos
		if !Q.CanAddVideo(req.RemoteAddr) {
			w.Header().Set("Content-type", "text/html")
			vidNotAddedTempl, _ := template.ParseFiles("templates/not_added.html")
			vidNotAddedTempl.Execute(w, nil)
			return
		}

		// Add video
		Q.QuickAddVideoLink(req.RemoteAddr, videoLink)

		w.Header().Set("Content-type", "text/html")
		vidAddedTempl, _ := template.ParseFiles("templates/added.html")
		vidAddedTempl.Execute(w, nil)
	} else {
		// Redirect back to homepage if not a POST request)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

// Endpoint for user to remove their own video from the queue
func userRemoveHandler(w http.ResponseWriter, req *http.Request) {
	// Get video id from post data
	if req.Method == http.MethodPost {
		id := req.PostFormValue("video_id")
		ip := help.GetIP(req.RemoteAddr)
		Q.UserRemoveVideo(id, ip)
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Endpoint for user to upload a file
func fileUploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		// Must open file before all else or connection is reset if redirected
		// Opens the POST'ed file
		file, header, err := req.FormFile("video_file")
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, "Can't parse uploaded file", err)
			return
		}

		_, aliasExists := Q.GetAlias(req.RemoteAddr)
		// If user has a valid alias
		if !aliasExists {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
			return
		}

		// If user has max added videos
		if !Q.CanAddVideo(req.RemoteAddr) {
			w.Header().Set("Content-type", "text/html")
			vidNotAddedTempl, _ := template.ParseFiles("templates/not_added.html")
			vidNotAddedTempl.Execute(w, nil)
			return
		}

		// New video id and destination file
		id := help.GenUUID()
		path := Q.DownloadFolder + "/" + id

		// Creates destination video file
		out, err := os.Create(path)
		defer out.Close()
		if err != nil {
			fmt.Fprintln(w, "Unable to create the video file for writing.")
			return
		}

		// write the video file to disk
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		// Add the video to queue
		go Q.AddUploadedVideo(req.RemoteAddr, header.Filename, path, id)

		// Return video added page
		w.Header().Set("Content-type", "text/html")
		vidAddedTempl, _ := template.ParseFiles("templates/added.html")
		vidAddedTempl.Execute(w, nil)
	} else {
		// Redirect back to homepage if not a POST request)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}