package youtube

// Simple wrapper of youtube-dl functions

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	"../config"
)

type Downloader struct {
	executable string
	downloadFolder string
}

func (ytdl *Downloader) Init() {
	ytdl.executable = config.Config.YTDLBin
	ytdl.downloadFolder = config.Config.DownloadFolder

	if !ytdl.checkFiles() {
		os.Mkdir(ytdl.downloadFolder, 0755)
		log.Println("Creating download folder:", ytdl.downloadFolder)
		if !ytdl.checkFiles() {
			log.Println("Error with folder creation:", ytdl.downloadFolder)
			os.Exit(1)
		}
	}
}

func (ytdl *Downloader) checkFiles() bool {
	if _, err := os.Stat(ytdl.executable); os.IsNotExist(err) {
		log.Println("Cannot find youtube-dl executable", ytdl.executable)
		os.Exit(1)
	}

	if _, err := os.Stat(ytdl.downloadFolder); os.IsNotExist(err) {
		log.Println("Download folder not found:", ytdl.downloadFolder)
		return false
	} else {
		log.Println("Download folder exists!")
		return true
	}
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
	outputPath := ytdl.downloadFolder + "/" + uuid
	// Download Video
	dl := exec.Command(ytdl.executable, "-o", outputPath, "--no-playlist", link)
	dl.Run()

	// Uses wildcard search for file extension of vid file with uuid name
	vidPath, _ := filepath.Glob((ytdl.downloadFolder+"/"+uuid+".*"))

	// Return first instance of file search xor error
	if len(vidPath) > 0 {
		return vidPath[0], false
	} else {
		return "", true
	}
}