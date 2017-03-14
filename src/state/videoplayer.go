package state

import (
	"log"
	"time"
	"os"
	"os/exec"
)

func (q *ProcessQueue) VideoPlayerService() {
	log.Println("Video player service start...")

	for {
		currentVid := q.GetNextVideo()

		// Set NowPlaying
		q.NPLock.Lock()
		q.NowPlaying = currentVid
		q.NPLock.Unlock()

		q.UpdateBucketCache()

		emptyVid := Video{}

		if currentVid != emptyVid {
			// Add video to played videos map
			q.ListLock.Lock()
			q.JustPlayed[currentVid.IpAddr] = currentVid
			q.ListLock.Unlock()
			// Explicit unlocks to prevent lock during entire video

			log.Println("Playing Video:", currentVid.Title)

			// Since Go is a bit weird here's some extra comments for how timeout is done
			// Make a message channel, size of one, and only transport errors
			timeoutChannel := make(chan error, 1)

			// Need to have whitespace splitted arguments as individual strings or as an array
			args := append(q.PlayerArgs, currentVid.File)

			// Set the player off and call wait in its own goroutine
			// It will send it's exit signal and/or errors to the timeoutChannel when done
			player := exec.Command(q.PlayerExe, args...)
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

				// The timoutChannel will unblock after the video player exits
				case err := <-timeoutChannel:
					if err != nil {
						log.Printf("Video player exited with error = %v", err)
					} else {
						log.Print("Video completed")
					}
			}

			// Delete played video file
			os.Remove(currentVid.File)
			log.Println("Removed video file:", currentVid.Title)

		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
