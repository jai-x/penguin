package state

// These data structures are derived from information in ProcessQueue
// These are used for user facing data requests

// Video struct to return to user
type VideoInfo struct {
	Title string
	IpAddr string
	Uploader string
	ID string
	Ready bool
}

// Type to return to client
type PlaylistInfo struct {
	Playlist [][]VideoInfo
	NowPlaying VideoInfo
	UserAlias string
}

// Method to convert a Video struct to VideoInfo struct, given aliasmap
// Essentially, updates using facing information from the playlist
func (q *ProcessQueue) ConvertToInfo(v Video) VideoInfo {
	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

	name, exists := q.Aliases[v.IpAddr]
	if !exists {
		name = "Anon"
	}
	return VideoInfo{v.Title, v.IpAddr, name, v.ID, v.Ready}
}

// Creates 2d slice to represent a bucket structure and stores it in BucketCache
func (q *ProcessQueue) UpdateBucketCache() {
	// Read lock for playlist
	q.ListLock.RLock()
	defer q.ListLock.RUnlock()

	// Write lock for cache
	q.CacheLock.Lock()
	defer q.CacheLock.Unlock()

	// clear the bucket to rebuild
	q.BucketCache = make([][]VideoInfo, q.buckets)

	// Temp map of ip to bucket it should be placed into
	ipToBucket := map[string]int{}

	for _, vid := range q.Playlist {
		// Get the bucket the video should be in, zero valued
		b, _ := ipToBucket[vid.IpAddr]

		// Convert video to VideoInfo and add to bucket
		q.BucketCache[b] = append(q.BucketCache[b], q.ConvertToInfo(vid))

		// Increment bucket for ip address
		ipToBucket[vid.IpAddr]++
	}
}

// Uses BucketCache and Aliasmap to generate user facing information
func (q *ProcessQueue) GetPlaylistInfo(addr string) PlaylistInfo {
	// New playlist information
	var out PlaylistInfo

	// Read lock the bucket cache
	q.CacheLock.RLock()
	defer q.CacheLock.RUnlock()
	// Use cache to populate the PlaylistInfo
	out.Playlist = q.BucketCache

	// Get the alias
	name, exists := q.GetAlias(addr)
	out.UserAlias = name
	if !exists {
		out.UserAlias = "Anon"
	}

	// Read lock NowPlaying
	q.NPLock.RLock()
	defer q.NPLock.RUnlock()
	// Get nowplaying
	out.NowPlaying = q.ConvertToInfo(q.NowPlaying)

	return out
}
