package musicserver

import (
	"log"
	"strings"
	"path/filepath"
	"net/http"
	"errors"
	"time"
	"os"
	"io"

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

func queueLink(req *http.Request) error {
	ip := getIPFromRequest(req)

	// Check if alias set for this ip
	alias, aliasSet := al.Alias(ip)
	if !aliasSet {
		return errors.New("No user alias set")
	}

	newLink := req.PostFormValue("video_link")
	if len(newLink) == 0 {
		return errors.New("No video link provided")
	}

	if !pl.Available(ip) {
		return errors.New("Video not added, user has max videos queued")
	}

	subs := false
	if req.PostFormValue("download_subs") == "on" {
		subs = true
	}

	dur, _ := time.ParseDuration(req.PostFormValue("vid_offset"))
	offset := int(dur.Seconds())

	newVid := playlist.NewVideo(ip, alias)
	newVid.Subs = subs
	newVid.Offset = offset
	pl.AddVideo(newVid)
	go downloadVideo(newLink, newVid.UUID)
	return nil
}

func queueUploadedVideo(req *http.Request) error {
	ip := getIPFromRequest(req)

	file, header, err := req.FormFile("video_file")
	if file == nil {
		return errors.New("No file uploaded")
	}
	defer file.Close()

	if err != nil {
		return errors.New("Can't parse uploaded file")
	}

	// Check if alias set for this ip
	alias, aliasSet := al.Alias(ip)
	if !aliasSet {
		return errors.New("No user alias set")
	}

	if !pl.Available(ip) {
		return errors.New("Video not added, user has max videos queued")
	}

	newVid := playlist.NewVideo(ip, alias)
	// Gen file path with filename as uuid and get file extension from header
	newPath := vidFolder + "/" + newVid.UUID  + fileExt(header.Filename)

	// Create file
	newFile, err := os.Create(newPath)
	defer newFile.Close()
	if err != nil {
		return errors.New("Unable to create the video file for writing: \n" + err.Error())
	}

	// Write file
	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}

	// Add information to Video struct
	newVid.Title = stripFileExt(header.Filename)
	newVid.File = newPath
	newVid.Ready = true

	pl.AddVideo(newVid)
	return nil
}
