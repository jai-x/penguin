package mserver

// Type to use for page templating
type PageInfo struct {
	Playlist []VideoInfo
	NowPlaying VideoInfo
	UserAlias string
}

// Type to return to client as JSON
type PlaylistInfo struct {
	Playlist []VideoInfo
	NowPlaying VideoInfo
}

func getPlaylistInfo() PlaylistInfo {
	// Lock playlist and aliases
	Q.ListLock.RLock()
	defer Q.ListLock.RUnlock()
	Q.AliasLock.RLock()
	defer Q.AliasLock.RUnlock()

	// New playlist information
	var out PlaylistInfo

	// Convert playlist Video structs to VideoInfo structs
	for _, vid := range Q.Playlist {
		// Append info slice
		out.Playlist = append(out.Playlist, vid.ConvertToInfo(Q.Aliases))
	}

	// Lock now playing
	Q.NPLock.RLock()
	defer Q.NPLock.RUnlock()

	// Added now playing info struct to page info
	out.NowPlaying = Q.NowPlaying.ConvertToInfo(Q.Aliases)

	return out
}