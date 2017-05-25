package musicserver

import (
	"log"
	"strings"
	"path/filepath"
	"net/http"

	"./playlist"
	"./youtube"
)

type playlistInfo struct {
	UserAlias string
	Playlist  [][]playlist.Video
}

func newPlaylistInfo(ip string) playlistInfo {
	alias, _ := al.Alias(ip)
	out := playlistInfo{alias, pl.Playlist()}
	return out
}

func getIPFromRequest(req *http.Request) string {
	// Check for x forward header
	xfIPList := req.Header.Get("X-Forwarded-For")
	// Split comma separated ip list and get the first ip
	firstXfIP := strings.TrimSpace(strings.Split(xfIPList, ",")[0])

	// First xforwarded ip is non-empty, so return it
	if len(firstXfIP) > 0 {
		return firstXfIP
	} else {
		// Strip port number from the req.RemoteAddr address
		i := strings.LastIndex(req.RemoteAddr, ":")
		return req.RemoteAddr[:i]
	}
}

func url(relative string) string {
	return serverDomain + relative
}

func downloadVideo(newLink, uuid string) {
	// New downloader
	dl := youtube.NewDownloader(newLink, uuid)

	// fetch and set title
	title, err := dl.Title()
	if err != nil {
		log.Println("Title fetch error:", newLink, err.Error())
		pl.RemoveVideo(uuid)
		return
	}
	pl.SetTitle(uuid, title)

	// Download and set video file
	filepath, err := dl.Filepath()
	if err != nil {
		log.Println("Download error:", title, err.Error())
		pl.RemoveVideo(uuid)
		return
	}
	pl.SetFile(uuid, filepath)
}

// Return only file extension including the dot
func fileExt(file string) string {
	return filepath.Ext(file)
}

// Remove file extension 
func stripFileExt(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
