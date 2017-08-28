package player

import (
	"log"
	"os/exec"
	"time"
	"strconv"

	"../playlist"
)

type VideoPlayer struct {
	Playing  bool
	// Anonymous field
	playerExe  string
	playerArgs []string
	timeout    time.Duration

	killChan chan bool
}

func NewVideoPlayer(t, exe string, args []string, vid playlist.Video) *VideoPlayer {
	newTimeout, _ := time.ParseDuration(t)

	// Start at specific time of Video has offset
	if vid.Offset > 0 {
		args = append(args, []string{"--start", strconv.Itoa(vid.Offset)}...)
	}

	// Append video file to end of arguments to play
	args = append(args, vid.File)

	out := VideoPlayer{
		false,
		exe,
		args,
		newTimeout,
		make(chan bool),
	}
	return &out
}

func (v *VideoPlayer) Play() {
	log.Println("Video player started")

	p := exec.Command(v.playerExe, v.playerArgs...)
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
			log.Fatalln(
				"Failed to kill video player:",
				err.Error(),
				v.playerExe,
				v.playerArgs,
			)
		}

	case <-time.After(v.timeout):
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
