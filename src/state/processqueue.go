package state

import (
	"log"
	"sync"
	"time"
	"strings"
	"path/filepath"

	"../help"
	"../config"
	"../youtube"
)

// Youtube downloader
var YTDL youtube.Downloader

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
	DownloadFolder string

	PlayerExe string
	PlayerArgs[] string
}

// Struct initialiser
func (q *ProcessQueue) Init() {
	// Init empty slices
	q.Playlist = make([]Video, 0)
	// Init empty map
	q.JustPlayed = make(map[string]Video)
	q.Aliases = make(map[string]string)

	// Initialise youtube-dl downloader,
	YTDL.Init()
	YTDL.Update()

	// Set timeout and max buckets and download folder
	q.timeout = time.Duration(config.Config.VideoTimeout)
	q.buckets = config.Config.MaxBuckets
	q.DownloadFolder = config.Config.DownloadFolder

	// Set video player and arguments
	q.PlayerExe = config.Config.VideoPlayer
	q.PlayerArgs = strings.Fields(config.Config.VideoPlayerArgs)

	// Init cache
	q.BucketCache = make([][]VideoInfo, q.buckets)
}

// Gets alias from string map given ip address key
func (q *ProcessQueue) GetAlias(addr string) (string, bool) {
	ip := help.GetIP(addr)

	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

	alias, exists := q.Aliases[ip]
	return alias, exists
}

// Sets alias in string map using ip address as key
func (q *ProcessQueue) SetAlias(addr, alias string) {
	ip := help.GetIP(addr)

	q.AliasLock.Lock()
	// Templates auto-escape strings
	q.Aliases[ip] = alias
	q.AliasLock.Unlock()

	q.UpdateBucketCache()
}

// Checks if given ip address can upload a video to playlist or has reached limit
func (q *ProcessQueue) CanAddVideo(addr string) bool {
	// get uploader ip
	ip := help.GetIP(addr)

	q.ListLock.RLock()
	defer q.ListLock.RUnlock()

	// iterates lists, returns false if queued videos by user ip >= max buckets
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

// Returns video popped from front of playlist or empty video
// Uses Justplayed map to create psudeo bucket behaviour from a simple list
func (q *ProcessQueue) GetNextVideo() Video {
	q.ListLock.Lock()

	// If list is empty or top video is not ready return empty video
	if len(q.Playlist) < 1 || !q.Playlist[0].Ready{
		log.Println("(/'-')/  No Playable Videos in Playlist \\('-'\\)")
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

// Add user uploaded video to queue
func (q *ProcessQueue) AddUploadedVideo(addr, title, filePath, id string) {
	ip := help.GetIP(addr)

	// Trim file extension from filename title, if any
	title = strings.TrimSuffix(title, filepath.Ext(title))

	// Create video struct and add to queue
	newVideo := Video{id, title, filePath, ip, true}
	q.ListLock.Lock()
	q.Playlist = append(q.Playlist, newVideo)
	q.ListLock.Unlock()
	log.Println("Added to playlist:", title)

	q.UpdateBucketCache()
}

// Add placeholder struct to queue and begin video downloader
func (q *ProcessQueue) QuickAddVideoLink(addr, link string) {
	// Gen new uuid and get uploaders ip
	newId := help.GenUUID()
	ip := help.GetIP(addr)

	// Add video struct with available information and set Ready: false
	newVideo := Video{newId, "", "", ip, false}
	q.ListLock.Lock()
	q.Playlist = append(q.Playlist, newVideo)
	q.ListLock.Unlock()
	log.Println("Added to playlist:", link)

	q.UpdateBucketCache()

	// Start downloading process in new goroutine
	go q.DownloadVideo(newId, link)
}

// Download video and fill filepath into given video id in playlist
func (q *ProcessQueue) DownloadVideo(vidId, link string) {
	log.Println("Starting download of video link:", link)
	title, ok := YTDL.GetTitle(link)
	if !ok {
		log.Println("ERROR: Cannot get video title of link:", link)
		log.Println("ERROR: Download aborted:", link)
		q.AdminRemoveVideo(vidId)
		return
	}

	// Add title to video struct in playlist and continue downloading
	q.ListLock.Lock()
	for i, vid := range q.Playlist {
		if vidId == vid.ID {
			q.Playlist[i].Title = title
			break
		}
	}
	q.ListLock.Unlock()
	q.UpdateBucketCache()

	// Download video with given uuid as filename
	vidFilePath, err := YTDL.GetVideo(vidId, link)
	if err {
		log.Println("ERROR: Download failed:", link, title)
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