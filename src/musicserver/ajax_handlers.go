package musicserver

import (
	"os"
	"io"
	"strings"
	"html/template"
	"encoding/json"
	"net/http"

	"../help"
)

// Endpoint for queueing video via ajax
func ajaxQueueHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := map[string]string{}

	if req.Method == http.MethodPost {
		_, aliasExists := Q.GetAlias(req.RemoteAddr)

		if !aliasExists {
			out["Message"] = "No user alias set"
			out["Type"] = "error"
			json.NewEncoder(w).Encode(out)
			return
		}

		videoLink := req.PostFormValue("video_link")

		// Submitted video link is blank
		if len(strings.TrimSpace(videoLink)) == 0 {
			out["Message"] = "No video link given"
			out["Type"] = "error"
			json.NewEncoder(w).Encode(out)
			return
		}

		// If user has max added videos
		if !Q.CanAddVideo(req.RemoteAddr) {
			out["Message"] = "Video not added, user has too many videos"
			out["Type"] = "warn"
			json.NewEncoder(w).Encode(out)
			return
		}

		// Start video downloader in new goroutine so 
		Q.QuickAddVideoLink(req.RemoteAddr, videoLink)

		out["Message"] = "Video added"
		out["Type"] = "success"
		json.NewEncoder(w).Encode(out)
	} else {
		out["Message"] = "Use POST method"
		out["Type"] = "error"
		json.NewEncoder(w).Encode(out)
	}
}

func ajaxPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
	templ, _ := template.ParseFiles("templates/playlist.html")
	templ.Execute(w, plInfo)
}

func ajaxAdminPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
	templ, _ := template.ParseFiles("templates/admin_playlist.html")
	templ.Execute(w, plInfo)
}

// Endpoint for user to upload a file
func ajaxFileUploadHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	response := map[string]string{}
	if req.Method == http.MethodPost {
		// Must open file before all else or connection is reset if redirected
		// Opens the POST'ed file
		file, header, err := req.FormFile("video_file")
		defer file.Close()
		if err != nil {
			response["Message"] = "Can't parse uploaded file"
			response["Type"] = "error"
			json.NewEncoder(w).Encode(response)
			return
		}

		_, aliasExists := Q.GetAlias(req.RemoteAddr)
		// If user has a valid alias
		if !aliasExists {
			response["Message"] = "No user alias set"
			response["Type"] = "error"
			json.NewEncoder(w).Encode(response)
			return
		}

		// If user has max added videos
		if !Q.CanAddVideo(req.RemoteAddr) {
			response["Message"] = "Video not added, user has too many videos"
			response["Type"] = "warn"
			json.NewEncoder(w).Encode(response)
			return
		}

		// New video id and destination file
		id := help.GenUUID()
		path := Q.DownloadFolder + "/" + id

		// Creates destination video file
		out, err := os.Create(path)
		defer out.Close()
		if err != nil {
			response["Message"] = "Unable to create the video file for writing."
			response["Type"] = "error"
			json.NewEncoder(w).Encode(response)
			return
		}

		// write the video file to disk
		_, err = io.Copy(out, file)
		if err != nil {
			response["Message"] = err.Error()
			response["Type"] = "error"
			json.NewEncoder(w).Encode(response)
			return
		}

		// Add the video to queue
		go Q.AddUploadedVideo(req.RemoteAddr, header.Filename, path, id)

		response["Message"] = "File uploaded"
		response["Type"] = "success"
		json.NewEncoder(w).Encode(response)
	} else {
		response["Message"] = "Use POST method"
		response["Type"] = "error"
		json.NewEncoder(w).Encode(response)
	}
}