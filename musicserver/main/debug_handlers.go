package server

import (
	"encoding/json"
	"net/http"

	"../help"
	"../playlist"
)

func debugListHandler(w http.ResponseWriter, req *http.Request) {
	list, _ := playlist.GetAllInfo()

	enc := json.NewEncoder(w)
	enc.Encode(list)
}

func debugNPHandler(w http.ResponseWriter, req *http.Request) {
	_, np := playlist.GetAllInfo()

	enc := json.NewEncoder(w)
	enc.Encode(np)
}

func debugIPHandler(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)

	alias, _ := playlist.GetAlias(req.RemoteAddr)

	msg := struct {
		RawAddress string
		IP         string
		Alias      string
	}{
		req.RemoteAddr,
		help.GetIP(req.RemoteAddr),
		alias,
	}
	enc.Encode(msg)
}
