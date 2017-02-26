package state

import (
	"log"
	"sync"
	"html"
	"time"

	"../help"
	"../youtube"
)

// Youtube downloader
var YTDL youtube.Downloader

// Debug
var debugMode bool

// Type to represent a queued video
type Video struct {
	ID string
	Title string
	File string
	IpAddr string
}

// Main data structure, holds entire state of the video system
type Queue struct {
	// Slice of video structs as the playlist
	Playlist []Video
	ListLock sync.RWMutex

	// Map of ip addresses to aliases
	Aliases map[string]string
	AliasLock sync.RWMutex

	// Currently playing video
	NowPlaying Video
	NPLock sync.RWMutex

	timeout time.Duration
}

// Constructor
func (q *Queue) Init(t time.Duration, debug bool) {
	q.Playlist = make([]Video, 0)
	q.Aliases = make(map[string]string)

	// Initialise youtube-dl downloader,
	YTDL.Init("youtube-dl/youtube-dl", "/tmp/")
	debugMode = debug

	if !debugMode {
		YTDL.Update()
	}

	// Set timeout
	q.timeout = t
}

// Returns video popped from front of playlist or blankVideo
func (q *Queue) GetNextVideo() Video {
	q.ListLock.Lock()
	defer q.ListLock.Unlock()

	var out Video
	if len(q.Playlist) > 0 {
		// Pop from front
		out, q.Playlist = q.Playlist[0], q.Playlist[1:]
	} else {
		out = Video{}
	}
	return out
}

func (q *Queue) GetAlias(addr string) (string, bool) {
	ip := help.GetIP(addr)

	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

	alias, exists := q.Aliases[ip]
	return alias, exists
}

func (q *Queue) SetAlias(addr, alias string) {
	ip := help.GetIP(addr)

	q.AliasLock.Lock()
	defer q.AliasLock.Unlock()

	// Alias is escaped before saving to map
	q.Aliases[ip] = html.EscapeString(alias)
}

// Downloads and adds video to playlist, will silently fail
func (q *Queue) DownloadAndAddVideo(addr, link string) {
	title, ok := YTDL.GetTitle(link)
	if !ok {
		log.Println("Failed download of video", link, "\nFrom address:", addr)
		return
	}

	// #### DEBUG ####
	// Add non file video struct to list and return
	if debugMode {
		newId := help.GenUUID()
		newVideo := Video{newId, title, "", help.GetIP(addr)}

		q.ListLock.Lock()
		defer q.ListLock.Unlock()

		// Append to playlist
		q.Playlist = append(q.Playlist, newVideo)

		return
	}
	
	newId := help.GenUUID()

	// Download video with given uuid as filename
	log.Println("Starting download:", title)
	vidFilePath := YTDL.GetVideo(newId, link)
	log.Print("Download complete:", title)

	// Escape title before saving to video struct
	title = html.EscapeString(title)

	// Create new video struct
	newVideo := Video{newId, title, vidFilePath, help.GetIP(addr)}

	q.ListLock.Lock()
	defer q.ListLock.Unlock()

	// Append to playlist
	q.Playlist = append(q.Playlist, newVideo)	
}