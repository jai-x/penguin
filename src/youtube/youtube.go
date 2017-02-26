package youtube

import (
	"log"
	"os/exec"
	"io/ioutil"
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
		title := string(output)
		return title, true
	}
}

// Downloads video and returns path to video file
func (ytdl *Downloader) GetVideo(uuid, link string) string {
	// Template will cause youtube-dl to give the video the uuid as filename
	// Will also download to specific folder
	template := ytdl.downloadFolder + uuid
	// Download Video
	dl := exec.Command(ytdl.executable, "-o", template, "--no-playlist", link)
	dl.Run()

	// Find the file extension using wildcard search in the download folder
	vidPath := ytdl.downloadFolder+findWildcardFilename(ytdl.downloadFolder, uuid+".*")

	// Return full path to video
	return vidPath
}

// Wacky wildcard search function
func findWildcardFilename(folder, pattern string) string {
	// Get slice of files in folder
	vidFiles, _ := ioutil.ReadDir(folder)

	// Iterate over files
	for _, vid := range vidFiles {
		if found, _ := filepath.Match(pattern, vid.Name()); found {
			return vid.Name()
		}
	}
	// Empty string for no match
	return ""
}