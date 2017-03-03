package state

// These data structures are derived from information in Queue
// These are used for user facing data requests

// Video struct with only user relevant fields
type VideoInfo struct {
	Title string
	IpAddr string
	Uploader string
	ID string
}

// Type to return to client
type PlaylistInfo struct {
	Playlist [][]VideoInfo
	NowPlaying VideoInfo
	UserAlias string
	Downloading []string
}

// Method to convert a Video struct to VideoInfo struct, given aliasmap
func (v *Video) ConvertToInfo(aliasMap map[string]string) VideoInfo {
	name, exists := aliasMap[v.IpAddr]
	if !exists {
		name = "Anon"
	}
	return VideoInfo{v.Title, v.IpAddr, name, v.ID}
}

func (q *Queue) UpdateBucketCache() {
	// Read lock for playlist and aliases
	q.ListLock.RLock()
	defer q.ListLock.RUnlock()
	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

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
		q.BucketCache[b] = append(q.BucketCache[b], vid.ConvertToInfo(q.Aliases))

		// Increment bucket for ip address
		ipToBucket[vid.IpAddr]++
	}
}

// Method to return PlaylistInfo from the Queue
func (q *Queue) GetPlaylistInfo(addr string) PlaylistInfo {
	// New playlist information
	var out PlaylistInfo

	// Read lock the bucket cache
	q.CacheLock.RLock()
	defer q.CacheLock.RUnlock()
	// Use cache to populate the PlaylistInfo
	out.Playlist = q.BucketCache
	out.Downloading = q.Downloading

	// Get the alias
	out.UserAlias, _ = q.GetAlias(addr)

	// Read lock NowPlaying and Aliases
	q.NPLock.RLock()
	defer q.NPLock.RUnlock()
	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()
	// Get nowplaying
	out.NowPlaying = q.NowPlaying.ConvertToInfo(q.Aliases)

	return out
}