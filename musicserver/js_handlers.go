package musicserver

import (
	"encoding/json"
	"net/http"
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

	if err := queueLink(req); err != nil {
		msg := AJAXMessage{err.Error(), "error"}
		out.Encode(msg)
		return
	}

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

	if err := queueUploadedVideo(req); err != nil {
		msg := AJAXMessage{err.Error(), "error"}
		out.Encode(msg)
		return
	}

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
