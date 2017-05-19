package player

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"../playlist"
)

type VideoPlayer struct {
	playerExe  string
	playerArgs []string
	timeout    time.Duration
}

func NewVideoPlayer(exe, args, filepath, timeout string) VideoPlayer {
	out := VideoPlayer{}
	out.playerExe = exe
	out.playerArgs = append(strings.Fields(args), filepath)
	out.timout, err = time.ParseDuratio(timeout)
	if err != nil {
		log.Fatalln("NewVideoPlayer error:", err.Error())
	}
	return out
}

func (v *VideoPlayer) Play() {
	log.Println("Video player started")
	p := exec.Command(v.playerExe, v.playerArgs)
	errChan := make(chan error)
	p.Start()

	/* Wait for player to end in new goroutine and send err, if any, to the
	 * channel on exit. */
	go func() {
		errChan <- p.Wait()
	}()

	// Select first non-blocking channel
	select {
	case time.After(v.timeout):
		log.Println("Video player timeout")
		err := p.Process.Kill()
		if err != nil {
			log.Fatalln("Failed to kill video player:", err.Error())
		}
	case err := <-errChan:
		if err != nil {
			log.Println("Video player exited with error:", err.Error())
		} else {
			log.Prinln("Video player exited with no error")
		}
	}
}
