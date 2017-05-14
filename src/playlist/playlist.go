package playlist

import (
	"log"
	"sync"

	"../config"
	"../help"
	"../youtubeDL"
)

var (
	bucketLock sync.RWMutex
	buckets    [][]Video
	nowPlaying Video

	bucketNumber int
)

// Init
func Init() {
	log.Println("Playlist init...")
	bucketNumber = config.Config.MaxBuckets
	nowPlaying = Video{}
	buckets = make([][]Video, bucketNumber)
	aliasMap = make(map[string]string)

	// Make list of lists
	for b := range buckets {
		buckets[b] = make([]Video, 0)
	}
}

/* Check if a new video can be added from the specified ip address.
 * Returns false if the number of videos form the ip is >= that the number of
 * buckets. Returns true otherwise. */
func CanAddVideo(addr string) bool {
	ip := help.GetIP(addr)

	bucketLock.RLock()
	defer bucketLock.RUnlock()

	num := 0

	for _, bucket := range buckets {
		for _, vid := range bucket {
			if vid.IpAddr == ip {
				num++
			}
		}
	}

	if num >= bucketNumber {
		return false
	} else {
		return true
	}
}

/* Removes a video from the playlist with the specified uuid. Will not affect
 * the playlist if a video doesnt exist. */
func RemoveVideo(uuid string) {
	bucketLock.Lock()
	defer bucketLock.Unlock()

	out, in := 0, 0
	foundIp := ""

	/* Bubble target video down list swapping with indexes of other videos
	 * uploaded by same ip */
	for b, bucket := range buckets {
		for v, vid := range bucket {
			if vid.IpAddr == foundIp {
				// Swap videos at index
				buckets[out][in], buckets[b][v] = buckets[b][v], buckets[out][in]
				// Store new index
				out, in = b, v
			}

			if vid.UUID == uuid {
				foundIp = vid.IpAddr
				out, in = b, v
			}
		}
	}

	// Delete video bubbled to bottom
	if foundIp != "" {
		log.Println("Removed video:", buckets[out][in].Title)
		stateChange()
		buckets[out][in].DeleteFile()
		buckets[out] = append(buckets[out][:in], buckets[out][in+1:]...)
	}
}

/* Returns true if the given address is the uploader of the video with the
 * specified uuid. If the video with the uuid does not exist, return false */
func AddrOwnsVideo(addr, uuid string) bool {
	ip := help.GetIP(addr)

	bucketLock.RLock()
	defer bucketLock.RUnlock()

	for _, bucket := range buckets {
		for _, vid := range bucket {
			if vid.IpAddr == ip && vid.UUID == uuid {
				return true
			}
		}
	}
	return false
}

/* Add empty video to playlist in teh first bucket that does not contain a video
 * from the same IP address. Then call the video downloader to download the
 * specified link to populate the video struct with the given uuid */
func AddVideoLink(addr, link string) {
	uuid := help.GenUUID()
	ip := help.GetIP(addr)

	aliasLock.RLock()
	alias, _ := aliasMap[ip]
	aliasLock.RUnlock()

	// Create new video with given ip, alias and uuid values
	newVid := Video{uuid, "", "", ip, alias, false, false}

	bucketLock.Lock()
	defer bucketLock.Unlock()

	for index, bucket := range buckets {
		contains := false
		for _, vid := range bucket {
			if vid.IpAddr == ip {
				contains = true
			}
		}

		// Append to first bucket that does not conatain video from the ip
		if !contains {
			buckets[index] = append(buckets[index], newVid)
			break
		}
	}
	log.Println("New link queued:", link)
	stateChange()
	go downloadVideo(uuid, link)
}

/* Download the title and video file for the specified link and populate the
 * struct of the video with the given uuid in the playlist */
func downloadVideo(uuid, link string) {
	title, ok := youtubeDL.GetTitle(link)
	if !ok {
		log.Println("Cannot get video title for link:", link)
		log.Println("Download aborted")
		RemoveVideo(uuid)
		stateChange()
		return
	}

	// Add title to video struct in playlist and continue downloading
	bucketLock.Lock()
	for b, bucket := range buckets {
		for v, vid := range bucket {
			if vid.UUID == uuid {
				buckets[b][v].Title = title
				goto titleAdded
			}
		}
	}

titleAdded:
	bucketLock.Unlock()
	stateChange()
	log.Println("Got video title:", title)

	// Download video with given uuid as filename
	vidFilePath, err := youtubeDL.GetVideo(uuid, link)
	if err {
		log.Println("Cannot get video file:", title)
		log.Println("Download aborted")
		RemoveVideo(uuid)
		stateChange()
		return
	}

	// Add new video to playlist
	bucketLock.Lock()
	for b, bucket := range buckets {
		for v, vid := range bucket {
			if vid.UUID == uuid {
				buckets[b][v].File = vidFilePath
				buckets[b][v].Ready = true
				goto videoAdded
			}
		}
	}

videoAdded:
	bucketLock.Unlock()
	stateChange()
	log.Println("Added new video:", title)
}

// Add a pre-filled video struct to the playlist
func AddVideoStruct(newVid Video) {

	bucketLock.Lock()
	defer bucketLock.Unlock()

	for index, bucket := range buckets {
		contains := false
		for _, vid := range bucket {
			if vid.IpAddr == newVid.IpAddr {
				contains = true
			}
		}

		// Append to first bucket that does not conatain video from the ip
		if !contains {
			buckets[index] = append(buckets[index], newVid)
			break
		}
	}
	stateChange()
	log.Println("Added new video:", newVid.Title)
}

// Return entire bucket and now playing video
func GetAllInfo() ([][]Video, Video) {
	bucketLock.RLock()
	defer bucketLock.RUnlock()

	return buckets, nowPlaying
}

// Advances the playlist to set nowplaying to the next video
func AdvancePlaylist() {
	bucketLock.Lock()

	// Empty playlist
	if len(buckets[0]) < 1 {
		bucketLock.Unlock()
		nowPlaying = Video{}
		return
	}

	// Find and return a ready and unplayed video
	for v, vid := range buckets[0] {
		if !vid.Played {
			bucketLock.Unlock()
			if vid.Ready {
				buckets[0][v].Played = true
				nowPlaying = vid
				stateChange()
				return
			} else {
				nowPlaying = Video{}
				return
			}
		}
	}

	// if reaching this point, all videos in top bucket are played

	// Delete all files for top bucket
	for _, vid := range buckets[0] {
		vid.DeleteFile()
	}

	// Shift buckets up and create new last bucket
	buckets = append(buckets[1:], make([]Video, 0))
	log.Println("Video bucket cycle")

	bucketLock.Unlock()
	stateChange()
	AdvancePlaylist()
}

func GetNP() Video {
	bucketLock.RLock()
	defer bucketLock.RUnlock()

	return nowPlaying
}
