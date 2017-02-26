package mserver

import (
	"log"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"
)

type Downloader struct {
	executable string
	downloadFolder string
}

func (ytdl Downloader) update() {
	log.Println("Updating youtube-dl ...")
	cmd := exec.Command(ytdl.executable, "--update")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("youtube-dl updated")
}

// Returns title as string and boolean if successful
func (ytdl Downloader) getTitle(link string) (string, bool) {
	cmd := exec.Command(ytdl.executable, "--get-title",  "--no-playlist", link)
	// Output() runs the command and produces output
	output, err := cmd.Output()
	if err != nil {
		return "", false
	} else {
		title := string(output)
		return title, true
	}
}

func (ytdl Downloader) getVideo(uuid, link string) {
	// Template will cause youtube-dl to give the video the uuid as filename
	// Will also download to specific folder
	template := ytdl.downloadFolder + uuid
	// Download Video
	dl := exec.Command(ytdl.executable, "-o", template, "--no-playlist", link)
	dl.Run()

	// youtube-dl always keeps file extensions, so this will strip them
	// mv := exec.Command("mv", ytdl.downloadFolder + uuid + ".* " + ytdl.downloadFolder + uuid)
	// err := mv.Run()
	// if err != nil {
	// 	log.Print(err)
	// }


	// FOR SOME EFFING REASON THE COMMENTED CODE ABOVE DOESNT WORK
	// SO THIS WILL HAVE TO MAKE DO INSTEAD
	// I GUESS THIS WAY IS TECHNICALLY MORE PORTABLE BUT IT DOESN'T MAKE ME ANY LESS ANGRY
	vidPath := ytdl.downloadFolder+findWildcardFilename(ytdl.downloadFolder, uuid+".*")
	os.Rename(vidPath, ytdl.downloadFolder+uuid)

}

func findWildcardFilename(folder, pattern string) string {
	vidFiles, _ := ioutil.ReadDir(folder)

	// FUCK YOU AND OUR READABILITY
	var out string
	for _, vid := range vidFiles {
		if found, _ := filepath.Match(pattern, vid.Name()); found {
			out = vid.Name()
		}
	}
	return out
}