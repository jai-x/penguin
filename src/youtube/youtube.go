package youtube

import (
	"log"
	"strings"
	"os/exec"
	"path/filepath"
)

type Downloader struct {
	executable string
	downloadFolder string
}

func (ytdl *Downloader) Init(exe, dir string) {
	ytdl.executable = exe
	ytdl.downloadFolder = dir
}

// Updates the binary
func (ytdl *Downloader) Update() {
	log.Println("Updating youtube-dl ...")
	cmd := exec.Command(ytdl.executable, "--update")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("youtube-dl updated")
}

// Returns title as string and boolean if successful
func (ytdl *Downloader) GetTitle(link string) (string, bool) {
	cmd := exec.Command(ytdl.executable, "--get-title",  "--no-playlist", link)
	// Output() runs the command and produces output
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return "", false
	} else {
		title := strings.TrimSpace(string(output))
		return title, true
	}
}

// Downloads video and returns path to video file, or error
func (ytdl *Downloader) GetVideo(uuid, link string) (string, bool) {
	// Template will cause youtube-dl to give the video the uuid as filename
	// Will also download to specific folder
	template := ytdl.downloadFolder + uuid
	// Download Video
	dl := exec.Command(ytdl.executable, "-o", template, "--no-playlist", link)
	dl.Run()

	// Uses wildcard search for file extension of vid file with uuid name
	vidPath, _ := filepath.Glob((ytdl.downloadFolder+"/"+uuid+".*"))

	// Return first instance of file search
	if len(vidPath) > 0 {
		return vidPath[0], false
	} else {
		return "", true
	}
}