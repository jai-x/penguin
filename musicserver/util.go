package musicserver

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	return conf.ServerDomain + relative
}

func downloadVideo(newLink string, newVid playlist.Video) {
	// New download settings
	st, err := youtube.NewSettings(conf.YTDLExe, conf.FFMPEGExe, conf.VidFolder)
	if err != nil {
		log.Fatalln("Settings for youtube-dl are incorrect:", err.Error())
	}
	// New downloader
	dl := youtube.NewDownloader(newLink, newVid.UUID, newVid.Subs, st)

	// fetch and set title
	title, err := dl.Title()
	if err != nil {
		log.Println("Title fetch error:", newLink, err.Error())
		pl.RemoveVideo(newVid.UUID)
		return
	}
	pl.SetTitle(newVid.UUID, title)

	// Download and set video file
	filepath, err := dl.Filepath()
	if err != nil {
		log.Println("Download error:", title, err.Error())
		pl.RemoveVideo(newVid.UUID)
		return
	}
	pl.SetFile(newVid.UUID, filepath)
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
	ldur, _:= time.ParseDuration(req.PostFormValue("vid_length"))
	if ldur == 0{
		ldur, _ = time.ParseDuration(conf.VidTimout)
	}
	offset := int(dur.Seconds())
	length := int(ldur.Seconds())

	newVid := playlist.NewVideo(ip, alias)
	newVid.Hash = newLink
	
	// Check if video has been uploaded before, if so, check for intersection with previously played sections
	if pl.R9kmode {
		playedSubset, ok := rm[newVid.Hash]
		if ok {
			beforeStart := 0
			beforeEnd := 0
			for _, i := range playedSubset{
				if i <= offset{
				beforeStart = beforeStart + 1
				}
				if i < offset+length{
				beforeEnd = beforeEnd + 1
				}
			}
			if (beforeEnd != beforeStart) || (beforeEnd % 2 == 1){
				return errors.New("Video not added, no reposts in r9k mode")
			}
		}
		rm[newVid.Hash] = append(rm[newVid.Hash], offset)
		rm[newVid.Hash] = append(rm[newVid.Hash], offset + length)
	}
	
	newVid.Subs = subs
	newVid.Offset = offset
	newVid.Length = length	
	pl.AddVideo(newVid)
	go downloadVideo(newLink, newVid)
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
	newPath := conf.VidFolder + "/" + newVid.UUID + fileExt(header.Filename)

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
	newVid.Hash = ip + newVid.Title
	playedSubset := rm[newVid.Hash]
	if playedSubset != nil {
		return errors.New("Video not added, no reposts in r9k mode")
	}

	pl.AddVideo(newVid)
	rm[newVid.Hash] = append(rm[newVid.Hash], 0)
	return nil
}
