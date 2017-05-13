package server

import (
	"os"
	"io"
	"net/http"
	"encoding/json"

	"../help"
	"../playlist"
)

type Message struct {
	Response string
	Type string
}

func ajaxQueueHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := json.NewEncoder(w)

	if req.Method != http.MethodPost {
		msg := Message{"Use POST Method", "error"}
		out.Encode(msg)
		return
	}

	_, aliasExists := playlist.GetAlias(req.RemoteAddr)
	if !aliasExists {
		msg := Message{"No user alias set", "error"}
		out.Encode(msg)
		return
	}

	newLink := req.PostFormValue("video_link")
	if len(newLink) == 0 {
		msg := Message{"No video link provided", "error"}
		out.Encode(msg)
		return
	}

	if !playlist.CanAddVideo(req.RemoteAddr) {
		msg := Message{"Video not added, user has max videos queued", "warn"}
		out.Encode(msg)
		return
	}

	playlist.AddVideoLink(req.RemoteAddr, newLink)
	msg := Message{"Video added", "success"}
	out.Encode(msg)
}

func ajaxUploadHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := json.NewEncoder(w)

	if req.Method != http.MethodPost {
		msg := Message{"Use POST Method", "error"}
		out.Encode(msg)
		return
	}

	file, header, err := req.FormFile("video_file")
	defer file.Close()
	if err != nil {
		msg := Message{"Can't parse uploaded file", "error"}
		out.Encode(msg)
		return
	}

	alias, aliasExists := playlist.GetAlias(req.RemoteAddr)
	if !aliasExists {
		msg := Message{"No user alias set", "error"}
		out.Encode(msg)
		return
	}

	if !playlist.CanAddVideo(req.RemoteAddr) {
		msg := Message{"Video not added, user has max videos queued", "warn"}
		out.Encode(msg)
		return
	}

	// Gen uuid and filepath
	uuid := help.GenUUID()
	path := vidFolder + "/" + uuid + help.GetFileExt(header.Filename)

	// Create file
	newFile, err := os.Create(path)
	defer newFile.Close()
	if err != nil {
		msg := Message{"Unable to create the video file for writing", "error"}
		out.Encode(msg)
		return
	}

	// Write file
	_, err = io.Copy(newFile, file)
	if err != nil {
		msg := Message{err.Error(), "error"}
		out.Encode(msg)
		return
	}

	// Struct for new video
	newVid := playlist.Video{
		UUID: uuid,
		Title: help.StripFileExt(header.Filename),
		File: path,
		IpAddr: help.GetIP(req.RemoteAddr),
		Alias: alias,
		Ready: true,
		Played: false,
	}

	playlist.AddVideoStruct(newVid)

	msg := Message{"File uploaded", "success"}
	out.Encode(msg)
}
