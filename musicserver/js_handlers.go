package musicserver

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"./playlist"
)

type AJAXMessage struct {
	Response string
	Type     string
}

func ajaxQueueHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := json.NewEncoder(w)
	ip := getIPFromRequest(req)

	if req.Method != http.MethodPost {
		msg := AJAXMessage{"Use POST Method", "error"}
		out.Encode(msg)
		return
	}

	alias, aliasExists := al.Alias(ip)
	if !aliasExists {
		msg := AJAXMessage{"No user alias set", "error"}
		out.Encode(msg)
		return
	}

	newLink := req.PostFormValue("video_link")
	if len(newLink) == 0 {
		msg := AJAXMessage{"No video link provided", "error"}
		out.Encode(msg)
		return
	}

	if !pl.Available(ip) {
		msg := AJAXMessage{"Video not added, user has max videos queued", "warn"}
		out.Encode(msg)
		return
	}

	subs := false
	if req.PostFormValue("download_subs") == "on" {
		subs = true
	}

	newVid := playlist.NewVideo(ip, alias, subs)
	pl.AddVideo(newVid)
	go downloadVideo(newLink, newVid.UUID)

	msg := AJAXMessage{"Video added", "success"}
	out.Encode(msg)
}

func ajaxUploadHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := json.NewEncoder(w)
	ip := getIPFromRequest(req)

	if req.Method != http.MethodPost {
		msg := AJAXMessage{"Use POST Method", "error"}
		out.Encode(msg)
		return
	}

	file, header, err := req.FormFile("video_file")
	if file == nil {
		msg := AJAXMessage{"No file uploaded", "error"}
		out.Encode(msg)
		return
	}
	defer file.Close()

	if err != nil {
		msg := AJAXMessage{"Can't parse uploaded file", "error"}
		out.Encode(msg)
		return
	}

	alias, aliasExists := al.Alias(ip)
	if !aliasExists {
		msg := AJAXMessage{"No user alias set", "error"}
		out.Encode(msg)
		return
	}

	if !pl.Available(ip) {
		msg := AJAXMessage{"Video not added, user has max videos queued", "warn"}
		out.Encode(msg)
		return
	}

	newVid := playlist.NewVideo(ip, alias, false)
	// Gen file path with filename as uuid and get file extension from header
	newPath := vidFolder + "/" + newVid.UUID  + fileExt(header.Filename)

	// Create file
	newFile, err := os.Create(newPath)
	defer newFile.Close()
	if err != nil {
		msg := AJAXMessage{"Unable to create the video file for writing", "error"}
		out.Encode(msg)
		return
	}

	// Write file
	_, err = io.Copy(newFile, file)
	if err != nil {
		msg := AJAXMessage{err.Error(), "error"}
		out.Encode(msg)
		return
	}

	// Add information to Video struct
	newVid.Title = stripFileExt(header.Filename)
	newVid.File = newPath
	newVid.Ready = true

	pl.AddVideo(newVid)

	msg := AJAXMessage{"File uploaded", "success"}
	out.Encode(msg)
}

func ajaxPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	ip := getIPFromRequest(req)
	info := newPlaylistInfo(ip)
	tl.Render(w, "playlist", info)
}

func ajaxAdminPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	ip := getIPFromRequest(req)
	info := newPlaylistInfo(ip)
	tl.Render(w, "admin_playlist", info)
}
