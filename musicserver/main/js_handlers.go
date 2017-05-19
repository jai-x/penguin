package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"../admin"
	"../help"
	"../playlist"
)

type AJAXMessage struct {
	Response string
	Type     string
}

func ajaxQueueHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := json.NewEncoder(w)

	if req.Method != http.MethodPost {
		msg := AJAXMessage{"Use POST Method", "error"}
		out.Encode(msg)
		return
	}

	_, aliasExists := playlist.GetAlias(req.RemoteAddr)
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

	if !playlist.CanAddVideo(req.RemoteAddr) {
		msg := AJAXMessage{"Video not added, user has max videos queued", "warn"}
		out.Encode(msg)
		return
	}

	playlist.AddVideoLink(req.RemoteAddr, newLink)
	msg := AJAXMessage{"Video added", "success"}
	out.Encode(msg)
}

func ajaxUploadHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := json.NewEncoder(w)

	if req.Method != http.MethodPost {
		msg := AJAXMessage{"Use POST Method", "error"}
		out.Encode(msg)
		return
	}

	file, header, err := req.FormFile("video_file")
	defer file.Close()
	if err != nil {
		msg := AJAXMessage{"Can't parse uploaded file", "error"}
		out.Encode(msg)
		return
	}

	alias, aliasExists := playlist.GetAlias(req.RemoteAddr)
	if !aliasExists {
		msg := AJAXMessage{"No user alias set", "error"}
		out.Encode(msg)
		return
	}

	if !playlist.CanAddVideo(req.RemoteAddr) {
		msg := AJAXMessage{"Video not added, user has max videos queued", "warn"}
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

	// Struct for new video
	newVid := playlist.Video{
		UUID:   uuid,
		Title:  help.StripFileExt(header.Filename),
		File:   path,
		IpAddr: help.GetIP(req.RemoteAddr),
		Alias:  alias,
		Ready:  true,
		Played: false,
	}

	playlist.AddVideoStruct(newVid)

	msg := AJAXMessage{"File uploaded", "success"}
	out.Encode(msg)
}

// helper function to send a string through a SSE message
func sseSend(w http.ResponseWriter, message string) {
	// Have to prefix all newlines with "data: " to send multline strings
	message = strings.Replace(message, "\n", "\ndata: ", -1)
	// Must prefix with "data: " and postfix with double newline
	fmt.Fprintf(w, "data: "+message+"\n\n")
}

func ssePlaylistHandler(w http.ResponseWriter, req *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Channel to notify if connection closes
	notify := w.(http.CloseNotifier).CloseNotify()
	// Channel for events in the playlist
	event := make(chan bool)

	// Set headers for SSE connection
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Check incoming ip for alias
	_, aliasExists := playlist.GetAlias(req.RemoteAddr)
	if !aliasExists {
		sseSend(w, "IP address has no alias set")
		return
	}

	connected := true
	for connected {
		go func() {
			playlist.WaitForChange()
			event <- true
		}()
		// Block until a channel returns
		select {
		// A playlist event has occurred
		case <-event:
			tmpl, err := template.ParseFiles("templates/playlist.html")
			if err != nil {
				log.Println(err.Error())
			}
			// Buffer to store rendered template
			var tmplBuffer bytes.Buffer
			// Execute template into buffer
			tmpl.Execute(&tmplBuffer, fetchInfo(req.RemoteAddr))
			// Send
			sseSend(w, tmplBuffer.String())
			f.Flush()
		// The connection has closed
		case <-notify:
			connected = false
		}
	}
}

func sseAdminPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Channel to notify if connection closes
	notify := w.(http.CloseNotifier).CloseNotify()
	// Channel for events in the playlist
	event := make(chan bool)

	// Set headers for SSE connection
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Check incoming ip for valid admin session
	if !admin.ValidSession(req.RemoteAddr) {
		sseSend(w, "IP address is not authenticated")
		return
	}

	connected := true
	for connected {
		go func() {
			playlist.WaitForChange()
			event <- true
		}()
		// Block until a channel returns
		select {
		// A playlist event has occurred
		case <-event:
			tmpl, err := template.ParseFiles("templates/admin_playlist.html")
			if err != nil {
				log.Println(err.Error())
			}
			// Buffer to store rendered template
			var tmplBuffer bytes.Buffer
			// Execute template into buffer
			tmpl.Execute(&tmplBuffer, fetchInfo(req.RemoteAddr))
			// Send
			sseSend(w, tmplBuffer.String())
			f.Flush()
		// The connection has closed
		case <-notify:
			connected = false
		}
	}
}
