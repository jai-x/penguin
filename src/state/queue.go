package state

import (
	"os"
	"log"
	"sync"
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
	// Video slice as the playlist
	Playlist []Video
	// Map of ip address to video that have been played
	JustPlayed map[string]Video
	// Mutex for both of the above structures
	ListLock sync.RWMutex

	// Map of ip addresses to aliases
	Aliases map[string]string
	AliasLock sync.RWMutex

	// Currently playing video
	NowPlaying Video
	NPLock sync.RWMutex

	// Cache to populate to PlaylistInfo structs
	BucketCache [][]VideoInfo
	Downloading []string
	CacheLock sync.RWMutex

	timeout time.Duration
	buckets int
}

// Struct initialiser
func (q *Queue) Init(t time.Duration, b int, debug bool) {
	debugMode = debug

	// Init empty slices
	q.Playlist = make([]Video, 0)
	// Init empty map
	q.JustPlayed = make(map[string]Video)
	q.Aliases = make(map[string]string)

	// Initialise youtube-dl downloader,
	YTDL.Init("youtube-dl/youtube-dl", "/tmp/")

	if !debugMode {
		YTDL.Update()
	}

	// Set timeout and max buckets
	q.timeout = t
	q.buckets = b

	// Init cache
	q.BucketCache = make([][]VideoInfo, q.buckets)
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
	// Templates auto-escape strings
	q.Aliases[ip] = alias
	q.AliasLock.Unlock()

	q.UpdateBucketCache()
}

func (q *Queue) CanAddVideo(addr string) bool {
	ip := help.GetIP(addr)

	q.ListLock.RLock()
	defer q.ListLock.RUnlock()

	ipQueue := 0
	for _, vid := range q.Playlist {
		if vid.IpAddr == ip {
			ipQueue++
			if ipQueue >= q.buckets {
				return false
			}
		}
	}
	return true
}

func (q *Queue) AddToDLing(title string) {
	q.CacheLock.Lock()
	q.Downloading = append(q.Downloading, title)
	q.CacheLock.Unlock()
}


func (q *Queue) RemoveFromDLing(title string) {
	q.CacheLock.Lock()
	for index, dl := range q.Downloading {
		if dl == title {
			q.Downloading = append(q.Downloading[:index], q.Downloading[index+1:]...)
			break
		}
	}
	q.CacheLock.Unlock()
}

// Returns video popped from front of playlist or blankVideo
func (q *Queue) GetNextVideo() Video {
	q.ListLock.Lock()

	// #### DEBUG CODE ####
	if len(q.Playlist) < 1 || debugMode {
		// return empty video if list is empty
		q.ListLock.Unlock()
		//log.Println("Playlist is empty")
		return Video{}
	}

	for index, vid := range q.Playlist {
		_, ipPlayed := q.JustPlayed[vid.IpAddr]
		// find the first video that does not have an ip in the JustPlayed map and in which the file exists
		if !ipPlayed {
			// Remove the video from slice at index
			q.Playlist = append(q.Playlist[:index], q.Playlist[index+1:]...)
			// Return the video
			q.ListLock.Unlock()
			return vid
		}
	}

	// If reaching this point the current bucket is empty
	q.JustPlayed = make(map[string]Video) // Delete map
	log.Println("Reached end of bucket")
	// Explicit unlock as this function will execute recursivly on new bucket
	q.ListLock.Unlock()
	// Recurse function
	return q.GetNextVideo()
}

// Download video and adds new video struct to playlist
func (q *Queue) DownloadAndAddVideo(addr, link string) {
	log.Println("Fetching title:", link)
	title, ok := YTDL.GetTitle(link)
	if !ok {
		log.Println("Failed download of video", link, "\nFrom address:", addr)
		return
	}

	newId := help.GenUUID()
	// #### DEBUG CODE ####
	// Add non file video struct to list and return
	if debugMode {
		newVideo := Video{newId, title, "", help.GetIP(addr)}
		// double check
		if q.CanAddVideo(addr) {
			q.ListLock.Lock()
			// Append to playlist
			q.Playlist = append(q.Playlist, newVideo)
			q.ListLock.Unlock()
			log.Println("DEBUG: Added video", title)
		} else {
			log.Println("DEBUG: Cannot add video, queue full", title)
		}

		log.Println("Updating bucket cache")
		q.UpdateBucketCache()
		return
	}

	log.Println("Starting download of new video:", title)
	// Add to downloading list
	q.AddToDLing(title)

	// Download video with given uuid as filename
	vidFilePath, err := YTDL.GetVideo(newId, link)
	if err {
		log.Println("Download failed:", title)
		q.RemoveFromDLing(title)
		return
	}
	log.Println("Download complete:", title)

	// Remove from downloading list
	q.RemoveFromDLing(title)

	// Doublecheck in case another video was added while this video was D/L'ing
	if q.CanAddVideo(addr) {
		// Add new video to playlist
		q.ListLock.Lock()
		newVideo := Video{newId, title, vidFilePath, help.GetIP(addr)}
		q.Playlist = append(q.Playlist, newVideo)
		q.ListLock.Unlock()
		log.Println("Video added:", title)
	} else {
		// remove unused video
		os.Remove(vidFilePath)
		log.Println("Cannot add video, too many in queue:", title)
	}

	log.Println("Updating bucket cache")
	q.UpdateBucketCache()
}

// This removes video and also bubbles the users video from the lower buckets upwards
func (q *Queue) AdminRemoveVideo(remVidId string) {
	q.ListLock.Lock()

	foundIP := ""
	prevIndex := 0

	// Iterate over playlist
	for index, vid := range q.Playlist {
		// This will evaluate to false until the ip is found at which point prevIndex is the video to remove
		// this will then bubble the target video down to lowest position that userip occupies
		if vid.IpAddr == foundIP {
			q.Playlist[prevIndex] = q.Playlist[index]
			prevIndex = index
		}

		// Find the video to remove via ID and get user ip
		if vid.ID == remVidId {
			prevIndex = index
			foundIP = vid.IpAddr
		}
	}

	// foundIp will be empty if given arguments are not valid
	if foundIP != "" {
		// prevIndex is now the target video in the lowest position so delete it
		q.Playlist = append(q.Playlist[:prevIndex], q.Playlist[prevIndex+1:]...)
	}

	q.ListLock.Unlock()
	q.UpdateBucketCache()
}

// Same as above function but requires both video id and video ip to sucessfully remove a video
func (q *Queue) UserRemoveVideo(remVidId, remUserIp string) {
	q.ListLock.Lock()

	foundIP := ""
	prevIndex := 0

	// Iterate over playlist
	for index, vid := range q.Playlist {
		// This will evaluate to false until the ip is found at which point prevIndex is the video to remove
		// this will then bubble the target video down to lowest position that userip occupies
		if vid.IpAddr == foundIP {
			q.Playlist[prevIndex] = q.Playlist[index]
			prevIndex = index
		}

		// Find the video to remove via ID and verify user trying to remove video owns it
		if vid.ID == remVidId && vid.IpAddr == remUserIp{
			prevIndex = index
			foundIP = vid.IpAddr
		}
	}

	// foundIp will be empty if given arguments are not valid
	if foundIP != "" {
		// prevIndex is now the target video in the lowest position so delete it
		q.Playlist = append(q.Playlist[:prevIndex], q.Playlist[prevIndex+1:]...)
	}

	q.ListLock.Unlock()
	q.UpdateBucketCache()
}