package playlist

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"
)

// Video must have valid UUID and IpAddr to be identifible
type Video struct {
	UUID   string
	Hash   string
	Title  string
	File   string
	IpAddr string
	Alias  string
	Offset int
	Length int
	Ready  bool
	Played bool
	NP     bool
	Subs   bool
}

// Creates an new Video struct with a pre-filled UUID variable.
func NewVideo(ip, alias string) Video {
	out := Video{}
	out.UUID = genUUID()
	out.IpAddr = ip
	out.Alias = alias
	out.Title = "New Video..."
	return out
}

// Get only the filename without the relative path
func (v *Video) RelativeFile() string {
	i := strings.LastIndex(v.File, "/")
	return v.File[i+1:]
}

// Delete the video file
func (v *Video) DeleteFile() {
	os.Remove(v.File)
}

// Generate a pseudorandom UUID string.
func genUUID() string {
	// Fill slice b with random bytes
	b := make([]byte, 16)
	_, err := rand.Read(b)
	// Error check
	if err != nil {
		log.Fatal("Cannot generate UUID: ", err.Error())
	}
	// Print to variable
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
