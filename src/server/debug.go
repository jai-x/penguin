package server

import (
	"net/http"
	"encoding/json"

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
