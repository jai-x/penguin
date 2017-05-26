package player

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"../playlist"
)

var (
	playerExe  string
	playerArgs []string
	timeout    time.Duration

)

func SetPlayer(exe, args string) {
	playerExe = exe
	// Splitting space sep string into list
	playerArgs = strings.Fields(args)
}

func SetTimeout(t string) error {
	err := error(nil)
	timeout, err = time.ParseDuration(t)
	return err
}

type VideoPlayer struct {
	Playing  bool
	Vid      playlist.Video

	killChan chan bool
}

func NewVideoPlayer(vid playlist.Video) VideoPlayer {
	out := VideoPlayer{}
	out.killChan = make(chan bool)
	out.Playing = false
	out.Vid = vid
	return out
}

func (v *VideoPlayer) Play() {
	log.Println("Video player started")

	// Set file to play by appending to end of arguments
	args := append(playerArgs, v.Vid.File)

	p := exec.Command(playerExe, args...)
	p.Start()
	v.Playing = true

	/* Wait for player to end in new goroutine and send err, if any, to the
	 * channel on exit. */
	errChan := make(chan error)
	go func() {
		errChan <- p.Wait()
	}()

	// Select first non-blocking channel
	select {
	case <-v.killChan:
		log.Println("Video player killed by admin")
		err := p.Process.Kill()
		if err != nil {
			log.Fatalln("Failed to kill video player:", err.Error())
		}

	case <-time.After(timeout):
		log.Println("Video player timeout")
		err := p.Process.Kill()
		if err != nil {
			log.Fatalln("Failed to kill video player:", err.Error())
		}

	case err := <-errChan:
		if err != nil {
			log.Println("Video player exited with error:", err.Error())
		} else {
			log.Println("Video player exited with no error")
		}
	}
	v.Playing = false
}

func (v *VideoPlayer) End() {
	if v.Playing {
		v.killChan <- true
	}
}
