package mserver

import (
	"fmt"
	"log"
	"sync"
	"html"
	"strings"
	"crypto/rand"
)

type Queue struct {
	Playlist []Video
	ListLock sync.RWMutex

	Aliases map[string]string
	AliasLock sync.RWMutex

	NowPlaying Video
	NPLock sync.RWMutex
}

func newQueue() Queue {
	q := Queue{}
	q.Playlist = make([]Video, 0)
	q.Aliases = make(map[string]string)
	return q
}

func (q *Queue) getNextVideo() Video {
	q.ListLock.Lock()
	defer q.ListLock.Unlock()

	var out Video
	if len(q.Playlist) > 0 {
		// Pop from front
		out, q.Playlist = q.Playlist[0], q.Playlist[1:]
	} else {
		out = BlankVideo
	}
	return out
}

func (q Queue) getAliasFromAddress(addr string) (string, bool) {
	ip := getIP(addr)

	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

	alias, exists := q.Aliases[ip]
	return alias, exists
}

func (q *Queue) setNewAlias(addr, alias string) {
	ip := getIP(addr)

	q.AliasLock.Lock()
	defer q.AliasLock.Unlock()

	// Alias is escaped before saving to map
	q.Aliases[ip] = html.EscapeString(alias)
}

// Downloads and adds video to playlist, will silently fail
func (q *Queue) downloadAndAddVideo(addr, link string) {
	title, ok := YTDL.getTitle(link)
	if !ok {
		log.Println("Failed download of video", link, "\nFrom address:", addr)
		return
	}
	
	newId := genUUID()

	// Download video with given uuid as filename
	log.Println("Starting download:", title)
	YTDL.getVideo(newId, link)
	log.Print("Download complete:", title)

	// Escape title before saving to video struct
	title = html.EscapeString(title)

	// Create new video struct
	newVideo := Video{title, YTDL.downloadFolder+newId, getIP(addr)}

	q.ListLock.Lock()
	defer q.ListLock.Unlock()

	// Append to playlist
	q.Playlist = append(q.Playlist, newVideo)	
}


// Addresses are in form "xxx.xxx.xxx.xxx:port"
// This strips the port number, returning only the IP
func getIP(addr string) string {
	return strings.Split(addr, ":")[0]
}

// Generate a pseudo random guid
// this is the only thing Go doesnt have in a standard lib
// TODO: Find a better way to do this
func genUUID() string {
	// Fill slice b with random bytes
	b := make([]byte, 16)
	_, err := rand.Read(b)
	// Error check
	if err != nil {
		log.Fatal("Error: ", err)
	}
	// Print to variable
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}