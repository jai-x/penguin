package state

import (
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
	Ready bool
}

// Main data structure, holds entire state of the video system
type ProcessQueue struct {
	// List lock mutex for both these data structures
	Playlist []Video
	JustPlayed map[string]Video
	ListLock sync.RWMutex

	// Map of ip addresses to aliases
	Aliases map[string]string
	AliasLock sync.RWMutex

	// Currently playing video
	NowPlaying Video
	NPLock sync.RWMutex

	// Cache to populate to PlaylistInfo structs
	BucketCache [][]VideoInfo
	CacheLock sync.RWMutex

	timeout time.Duration
	buckets int
}

// Struct initialiser
func (q *ProcessQueue) Init(t time.Duration, b int, debug bool) {
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

func (q *ProcessQueue) GetAlias(addr string) (string, bool) {
	ip := help.GetIP(addr)

	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

	alias, exists := q.Aliases[ip]
	return alias, exists
}

func (q *ProcessQueue) SetAlias(addr, alias string) {
	ip := help.GetIP(addr)

	q.AliasLock.Lock()
	// Templates auto-escape strings
	q.Aliases[ip] = alias
	q.AliasLock.Unlock()

	q.UpdateBucketCache()
}

func (q *ProcessQueue) CanAddVideo(addr string) bool {
	ip := help.GetIP(addr)

	q.ListLock.RLock()
	defer q.ListLock.RUnlock()

	ipProcessQueue := 0
	for _, vid := range q.Playlist {
		if vid.IpAddr == ip {
			ipProcessQueue++
			if ipProcessQueue >= q.buckets {
				return false
			}
		}
	}
	return true
}

// Returns video popped from front of playlist or blankVideo
func (q *ProcessQueue) GetNextVideo() Video {
	q.ListLock.Lock()

	// #### DEBUG CODE ####
	if len(q.Playlist) < 1 || debugMode {
		// return empty video if list is empty
		q.ListLock.Unlock()
		return Video{}
	}

	if !q.Playlist[0].Ready {
		// Top video is not ready so return empty video
		q.ListLock.Unlock()
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

func (q *ProcessQueue) QuickAddVideo(addr, link string) {
	newId := help.GenUUID()
	ip := help.GetIP(addr)

	newVideo := Video{newId,"", "", ip, false}
	q.ListLock.Lock()
	q.Playlist = append(q.Playlist, newVideo)
	q.ListLock.Unlock()

	q.UpdateBucketCache()

	go q.DownloadVideo(newId, link)
}

// Download video and updates video struct playlist
func (q *ProcessQueue) DownloadVideo(vidId, link string) {
	log.Println("Fetching title:", link)
	title, ok := YTDL.GetTitle(link)
	if !ok {
		log.Println("Failed stat video title of link:", link)
		q.AdminRemoveVideo(vidId)
		return
	}

	// Add title to video struct in playlist and continue downloading
	q.ListLock.Lock()
	for i, vid := range q.Playlist {
		if vidId == vid.ID {
			log.Println("Video found in list, adding title:", title)
			q.Playlist[i].Title = title
			if debugMode { q.Playlist[i].Ready = true }
			break
		}
	}
	q.ListLock.Unlock()
	q.UpdateBucketCache()

	// #### DEBUG ####
	if debugMode {
		log.Println("DEBUG: Skipping file download:", title)
		return
	}

	// Download video with given uuid as filename
	vidFilePath, err := YTDL.GetVideo(vidId, link)
	if err {
		log.Println("Download failed:", title)
		q.AdminRemoveVideo(vidId)
		return
	}

	// Add new video to playlist
	q.ListLock.Lock()
	for i, vid := range q.Playlist {
		if vid.ID == vidId {
			q.Playlist[i].File = vidFilePath
			q.Playlist[i].Ready = true
			break
		}
	}
	q.ListLock.Unlock()
	log.Println("Download completed:", title)

	q.UpdateBucketCache()
}

// This removes video and also bubbles the users video from the lower buckets upwards
func (q *ProcessQueue) AdminRemoveVideo(remVidId string) {
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
func (q *ProcessQueue) UserRemoveVideo(remVidId, remUserIp string) {
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

func (q *ProcessQueue) getAllAliases() map[string]string {
	q.AliasLock.RLock()
	out := q.Aliases
	q.AliasLock.RUnlock()
	return out
}
