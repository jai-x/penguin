package state

import (
	"log"
	"time"
	"os/exec"
)

func (q *Queue) PlayVideos() {
	for {
		currentVid := q.GetNextVideo()

		q.NPLock.Lock()
		q.NowPlaying = currentVid
		q.NPLock.Unlock()

		emptyVid := Video{}

		if currentVid != emptyVid {

			log.Println("Playing Video:", currentVid.Title)

			// Since Go is a bit weird here's some extra comments for how timeout is done
			// Make a message channel, size of one, and only transport errors
			timeoutChannel := make(chan error, 1)

			// Set the player off and call wait in its own goroutine
			// It will send its exit signal and/or errors to the timeoutChannel when done
			player := exec.Command("mpv", "-vo", "xv", "-fs", "-quiet", "-af=drc=2:0.25", currentVid.File)
			player.Start()
			go func() {
				timeoutChannel <- player.Wait()
			}()

			// The select switch will choose the first non-blocking channel
			select {
				// This empty channel will unblock after timer ends
				case <-time.After(q.timeout * time.Second):
					if err := player.Process.Kill(); err != nil {
						log.Fatal("Failed to kill video player:", err)
					}
					log.Println("Video reached timeout")

				// The timoutChannel will unblock after the video plyer exits
				case err := <-timeoutChannel:
					if err != nil {
						log.Printf("Video player exited with error = %v", err)
					} else {
						log.Print("Video completed")
					}
			}
		} else {
			log.Println("(/'-')/  No Videos in Playlist \\('-'\\)")
		}
		time.Sleep(1 * time.Second)
	}
}