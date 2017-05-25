package musicserver

import (
	"encoding/json"
	"net/http"
)

// Show the entire playlist as JSON
func debugListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	ip := getIPFromRequest(req)
	info := newPlaylistInfo(ip)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t") // To pretty print
	enc.Encode(info)
}

// Show information that the program has on the user IP address
func debugIPHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	ip := getIPFromRequest(req)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t") // To pretty print
	alias, _ := al.Alias(ip)
	// Anonymous struct to contain the information
	msg := struct {
		RawAddress string
		IP         string
		Alias      string
	}{
		req.RemoteAddr,
		ip,
		alias,
	}
	enc.Encode(msg)
}

func debugHeaderHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t") // To pretty print
	enc.Encode(req.Header)
}
