package videoplayer

import (
	"log"
	"time"
	"strings"
	"os/exec"

	"../config"
	"../playlist"
)

var (
	playerExe string
	playerArgs []string
	timeout time.Duration
)

func Init() {
	log.Println("Videoplayer init...")
	playerExe = config.Config.VideoPlayer
	// Arguments must be a slice of strings, whitespace separated
	playerArgs = strings.Fields(config.Config.VideoPlayerArgs)
	timeout, _= time.ParseDuration(config.Config.VideoTimeout)
}


func Start() {
	log.Println("Video player service start")

	emptyVid := playlist.Video{}

	for {
		playlist.AdvancePlaylist()
		newVid := playlist.GetNP()

		if newVid != emptyVid {
			// Channel to signal player exit or timout
			timeoutChan := make(chan error, 1)
			// Append file path to end or player arguments
			args := append(playerArgs, newVid.File)
			// Create player
			player := exec.Command(playerExe, args...)
			log.Println("Playing video:", newVid.Title)
			// Start the player
			player.Start()
			/* Wait for player to end in new goroutine and send message through
			 * channel on exit */
			go func() {
				timeoutChan <- player.Wait()
			}()

			// Select will choose the first non-blocking channel
			select {
			// This channel will unblock after timeout and kill the player
			case <-time.After(timeout):
				log.Println("Timeout reached, killing player")
				err := player.Process.Kill()
				if err != nil {
					log.Fatal("Failed to kill video player:", err.Error)
				}
				log.Println("Video reached timeout")
			// This channel will unblock when the player exits iteself
			case err := <-timeoutChan:
				if err != nil {
					log.Println("Video player exited with error:", err.Error)
				} else {
					log.Println("Video completed")
				}
			}
		} else {
			// Video is empty so wait and poll again later
			time.Sleep(2 * time.Second)
			log.Println("(/'-')/ No Videos to Play \\('-'\\)")
		}
	}
}
