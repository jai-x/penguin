package state

import "../help"

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
	Playlist []VideoInfo
	NowPlaying VideoInfo
	UserAlias string
}

// Method to convert a Video struct to VideoInfo struct, given aliasmap
func (v *Video) ConvertToInfo(aliasMap map[string]string) VideoInfo {
	name, exists := aliasMap[v.IpAddr]
	if exists {
		return VideoInfo{v.Title, v.IpAddr, name, v.ID}
	} else {
		return VideoInfo{v.Title, v.IpAddr, "Anon", v.ID}
	}
}

// Method to return PlaylistInfo from the Queue
func (q *Queue) GetPlaylistInfo(addr string) PlaylistInfo {
	// New playlist information
	var out PlaylistInfo

	// Lock playlist and aliases
	q.ListLock.RLock()
	defer q.ListLock.RUnlock()
	q.AliasLock.RLock()
	defer q.AliasLock.RUnlock()

	// Get alias
	alias, ok := q.Aliases[help.GetIP(addr)]
	if !ok {
		alias = "Anon"
	}
	out.UserAlias = alias

	// Convert playlist Video structs to VideoInfo structs
	for _, vid := range q.Playlist {
		// Append info slice
		out.Playlist = append(out.Playlist, vid.ConvertToInfo(q.Aliases))
	}

	// Lock nowplaying
	q.NPLock.RLock()
	defer q.NPLock.RUnlock()

	// Add now playing info struct
	out.NowPlaying = q.NowPlaying.ConvertToInfo(q.Aliases)

	return out
}