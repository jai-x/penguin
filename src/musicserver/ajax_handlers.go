package musicserver

import (
	"strings"
	"encoding/json"
	"net/http"
)

// Endpoint to return playlist JSON
func playlistHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	info := Q.GetPlaylistInfo(req.RemoteAddr)
	json.NewEncoder(w).Encode(info)
}

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
		go Q.DownloadAndAddVideo(req.RemoteAddr, videoLink)

		out["Message"] = "Video added"
		out["Type"] = "success"
		json.NewEncoder(w).Encode(out)
	} else {
		out["Message"] = "Use POST method"
		out["Type"] = "error"
		json.NewEncoder(w).Encode(out)
	}
}