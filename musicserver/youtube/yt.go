package youtube

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// Default use set variables

	// Path to youtube-dl executable
	ytExe string = "youtube-dl"
	// Path to ffmpeg executable, optionally set
	ffmpegExe string = ""
	// Path to folder to store downloaded videos
	dlPath string = "/tmp"
)

type Downloader struct {
	link string
	uuid string
}

func NewDownloader(link, uuid string) Downloader {
	// TODO: youtube link parse
	out := Downloader{}
	out.link = link
	out.uuid = uuid
	return out
}

// Sets the ytExe variable. If the file doesn't exist, will return error.
func SetYTDLPath(in string) error {
	if _, err := os.Stat(in); os.IsNotExist(err) {
		return err
	}
	ytExe = in
	return nil
}

// Sets the optional ffmpegExe variable. If file doesn't exist, will return
// error.
func SetFFMPEGPath(in string) error {
	if _, err := os.Stat(in); os.IsNotExist(err) {
		return err
	}
	ffmpegExe = in
	return nil
}

// Sets the dlPath variable. Will create the folder if it doesn't already
// exist. If there is an issue in folder creation, will return an error.
func SetDownloadFolder(in string) error {
	if _, err := os.Stat(in); os.IsNotExist(err) {
		// Folder doesn't exist so create it
		err = os.MkdirAll(in, 0755)
		if err != nil {
			// Return if there is an error creating the folder
			return err
		}
	}
	dlPath = in
	return nil
}

// Returns title of video of the given link.
func (d *Downloader) Title() (string, error) {
	args := []string{"--get-title", "--no-playlist", d.link}
	dl := exec.Command(ytExe, args...)
	output, err := dl.Output()
	title := strings.TrimSpace(string(output))
	return title, err
}

// Downloads video and returns filepath as string to downloaded video file.
func (d *Downloader) Filepath() (string, error) {
	// Use output template to download to folder and set filename as UUID
	outputPath := dlPath + "/" + d.uuid + `.%(ext)s`
	args := []string{"-o", outputPath, "--no-playlist"}

	// Optionally apply ffmpeg argument if variable is not empty
	if ffmpegExe != "" {
		args = append(args, []string{"--ffmpeg-location", ffmpegExe}...)
	}

	// Add link
	args = append(args, d.link)

	dl := exec.Command(ytExe, args...)
	err := dl.Run()
	if err != nil {
		// Return if any errors in downloading video
		return "", err
	}

	// The file extension for the new video is unknown, so perform a search of
	// the filename using a wildcard for the extension, if any.
	res, _ := filepath.Glob(dlPath + "/" + d.uuid +  "*")

	// No search results
	if len(res) < 1 {
		return "", errors.New("Cannot find downloaded video file: " + d.uuid)
	}

	// Return first search result
	return res[0], nil
}
