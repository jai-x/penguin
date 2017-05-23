package player

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

type VideoPlayer struct {
	playerExe    string
	playerArgs   []string
	timeout      time.Duration
	Playing      bool
	killChan     chan bool
}

func NewVideoPlayer(exe, args, filepath, timeout string) VideoPlayer {
	out := VideoPlayer{}
	out.killChan = make(chan bool)
	out.Playing = false
	out.playerExe = exe
	out.playerArgs = append(strings.Fields(args), filepath)
	var err error
	out.timeout, err = time.ParseDuration(timeout)
	if err != nil {
		log.Fatalln("NewVideoPlayer error:", err.Error())
	}
	return out
}

func (v *VideoPlayer) Play() {
	log.Println("Video player started")
	p := exec.Command(v.playerExe, v.playerArgs...)
	errChan := make(chan error)
	p.Start()
	v.Playing = true

	/* Wait for player to end in new goroutine and send err, if any, to the
	 * channel on exit. */
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
		v.Playing = false

	case <-time.After(v.timeout):
		log.Println("Video player timeout")
		err := p.Process.Kill()
		if err != nil {
			log.Fatalln("Failed to kill video player:", err.Error())
		}
		v.Playing = false

	case err := <-errChan:
		if err != nil {
			log.Println("Video player exited with error:", err.Error())
		} else {
			log.Println("Video player exited with no error")
		}
		v.Playing = false
	}
}

func (v *VideoPlayer) End() {
	if v.Playing {
		v.killChan <- true
	}
}
