package musicserver

import (
	"log"
	"strings"
	"path/filepath"

	"./playlist"
	"./youtube"
)

type playlistInfo struct {
	UserAlias string
	Playlist  [][]playlist.Video
}

func newPlaylistInfo(addr string) playlistInfo {
	ip := ip(addr)
	alias, _ := al.Alias(ip)
	out := playlistInfo{alias, pl.Playlist()}
	return out
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

// Addresses are in form "xxx.xxx.xxx.xxx:port"
// This strips the port number, returning only the IP
func ip(addr string) string {
	index := strings.LastIndex(addr, ":")
	ip := addr[:index]
	if ip == "[::1]" {
		return "localhost"
	}
	return ip
}

// Return only file extension including the dot
func fileExt(file string) string {
	return filepath.Ext(file)
}

// Remove file extension 
func stripFileExt(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
