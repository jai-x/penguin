package mserver

import (
	"log"
	"time"
	"os/exec"
	"net/http"
)

var (
	videoTimeout int

	YTDL Downloader

	BlankVideo Video = Video{"Nothing playing", "", ""}

	Q Queue
)


func Run() {

	YTDL = Downloader{"youtube-dl/youtube-dl", "/tmp/"}
	// Update youtube-dl
	//YTDL.update()

	// Set video timeout (in seconds)
	videoTimeout = 547

	// Initialize map and slice of global queue
	Q = newQueue()

	// Url Handlers
	http.HandleFunc("/alias", aliasHandler)
	http.HandleFunc("/queue", queueHandler)
	http.HandleFunc("/home", homeHandler)

	// Start video player function in a separate goroutine
	go playVideos()

	// Run the server
	log.Println("Running music server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func playVideos() {
	for {
		currentVid := Q.getNextVideo()

		Q.NPLock.Lock()
		Q.NowPlaying = currentVid
		Q.NPLock.Unlock()

		if currentVid != BlankVideo {

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
				case <-time.After(videoTimeout * time.Second):
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