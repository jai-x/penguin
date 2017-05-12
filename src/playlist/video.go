package playlist

import (
	"os"
	"strings"
)

type Video struct {
	UUID string
	Title string
	File string
	IpAddr string
	Alias string
	Ready bool
	Played bool
}

func (v *Video) DeleteFile() {
	os.Remove(v.File)
}

// returns only the filename and extension without the full path
func (v *Video) RelativeFile() string {
	ind := strings.LastIndex(v.File, "/")
	return v.File[ind+1:]
}
