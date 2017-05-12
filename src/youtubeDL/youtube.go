package youtubeDL

// Simple wrapper of youtube-dl functions

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	"../config"
	"../help"
)

var (
	executable string
	downloadFolder string
	ffmpegExe string
)

func Init() {
	log.Println("YoutubeDL init...")
	executable = config.Config.YTDLBin
	downloadFolder = config.Config.DownloadFolder
	ffmpegExe = config.Config.FFMPEGBin

	if !checkFiles() {
		os.Mkdir(downloadFolder, 0755)
		log.Println("Creating download folder:", downloadFolder)
		if !checkFiles() {
			log.Println("Error with folder creation:", downloadFolder)
			os.Exit(1)
		}
	}
}

func checkFiles() bool {
	if _, err := os.Stat(executable); os.IsNotExist(err) {
		log.Println("Cannot find youtube-dl executable", executable)
		os.Exit(1)
	}

	if _, err := os.Stat(downloadFolder); os.IsNotExist(err) {
		log.Println("Download folder not found:", downloadFolder)
		return false
	} else {
		log.Println("Download folder exists!")
		return true
	}
}

// Updates the binary
func Update() {
	log.Println("Updating youtube-dl ...")
	cmd := exec.Command(executable, "--update")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("youtube-dl updated")
}

// Returns title as string and boolean if successful
func GetTitle(link string) (string, bool) {
	cmd := exec.Command(executable, "--get-title", "--no-playlist", link)
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
func GetVideo(uuid, link string) (string, bool) {
	// Strip playlist info from URL if it is a youtube link
	link = help.StripYoutubePlaylist(link)

	// Template will cause youtube-dl to give the video the uuid as filename
	// Will also download to specific folder
	outputPath := downloadFolder + "/" + uuid
	// Download Video
	dl := exec.Command(executable, "-o", outputPath, "--ffmpeg-location", ffmpegExe, "--no-playlist", link)
	dl.Run()

	// Uses wildcard search for file extension of vid file with uuid name
	vidPath, _ := filepath.Glob((downloadFolder+"/"+uuid+".*"))

	// Return first instance of file search xor error
	if len(vidPath) > 0 {
		return vidPath[0], false
	} else {
		return "", true
	}
}
