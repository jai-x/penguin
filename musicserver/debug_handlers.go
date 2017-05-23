package musicserver

import (
	"encoding/json"
	"net/http"
)

// Show the entire playlist as JSON
func debugListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	info := newPlaylistInfo(req.RemoteAddr)
	enc := json.NewEncoder(w)
	// To pretty print
	enc.SetIndent("", "\t")
	enc.Encode(info)
}

// Show information that the program has on the user IP address
func debugIPHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	ip := ip(req.RemoteAddr)

	enc := json.NewEncoder(w)
	// To pretty print
	enc.SetIndent("", "\t")

	alias, _ := al.Alias(ip)

	// Anonymouse struct to contain the information
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
