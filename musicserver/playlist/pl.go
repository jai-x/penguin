package playlist

import (
	"sync"
)

type Playlist struct {
	mu       sync.RWMutex
	playlist [][]Video

	sublistNo int
	R9kmode   bool
}

func NewPlaylist(b int, c bool) Playlist {
	// Default bucket value is 4
	if b < 1 {
		b = 4
	}

	out := Playlist{}
	out.sublistNo = b
	out.R9kmode = c
	out.playlist = make([][]Video, out.sublistNo)

	// Init each subslice of Video in the playlist
	for index := range out.playlist {
		out.playlist[index] = make([]Video, 0)
	}

	return out
}

// Will alter change the number of sublists to the given value. If the new 
// sublist value is smaller than the current number of sublist, the trailing 
// sublists and their Videos will be deleted
func (p *Playlist) SetSublistCount(b int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if b < 1 {
		b = 4 // Default
	}

	// Make new playlist
	newPlaylist := make([][]Video, b)

	// Copy from existing playlist or create new sublist 
	for i := range newPlaylist {
		if i < len(p.playlist) {
			newPlaylist[i] = p.playlist[i]
		} else {
			newPlaylist[i] = make([]Video, 0)
		}
	}
	// Set 
	p.playlist = newPlaylist
}

// Returns true or false if the Playlist is available to take a new Video
// struct from the given IP address.
func (p *Playlist) Available(ip string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Simply check if an IP own a video in the last sublist
	last := len(p.playlist) - 1
	for _, vid := range p.playlist[last] {
		if vid.IpAddr == ip {
			return false
		}
	}
	return true
}

// Adds a new video struct to the playlist. Will append to the first sublist
// that does not conatain a Video with the same IP address. If all sublists
// contain a Video from the same IP address, the new Video will not be
// added.
func (p *Playlist) AddVideo(newVid Video) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for s, subl := range p.playlist {
		/* Find first sublist that does not contain a Video from the same IP
		 * address as the new Video */
		contains := false
		for _, vid := range subl {
			if vid.IpAddr == newVid.IpAddr {
				contains = true
			}
		}

		// Append the new Video to the first sublist that does not already
		// contain a Video from the same IP address and return.
		if !contains {
			p.playlist[s] = append(p.playlist[s], newVid)
			return
		}
	}
}

// Removes a Video struct from the playlist with a matching provided UUID
// string.
func (p *Playlist) RemoveVideo(remUUID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Index positions of last video with same IP address
	var out, in int
	// IP address of video to remove
	var remIP string

	for s, subl := range p.playlist {
		for v, vid := range subl {
			// Mark postion an IP address of Video to remove
			if vid.UUID == remUUID {
				out, in, remIP = s, v, vid.IpAddr
			}

			// Swap video down to lowest point occupied by the same IP address
			if vid.IpAddr == remIP {
				p.playlist[s][v], p.playlist[out][in] = p.playlist[out][in], p.playlist[s][v]
				out, in = s, v
			}
		}
	}

	// Delete video file
	p.playlist[out][in].DeleteFile()
	// Slice trick to remove Video struct from list while preserving sublist order.
	p.playlist[out] = append(p.playlist[out][:in], p.playlist[out][in+1:]...)
}

// Sets Title variable for a Video struct in the playlist with the matching
// UUID.
func (p *Playlist) SetTitle(vidUUID, newTitle string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for s, subl := range p.playlist {
		for v, vid := range subl {
			if vid.UUID == vidUUID {
				p.playlist[s][v].Title = newTitle
			}
		}
	}
}

//Sets Hash variable for a Video struct in the playlist with the matching UUID
func (p *Playlist) SetHash(vidUUID, newHash string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for s, subl := range p.playlist {
		for v, vid := range subl {
			if vid.UUID == vidUUID {
				p.playlist[s][v].Hash = newHash
			}
		}
	}
}

// Sets File variable for a Video struct in the playlist with the matching
// UUID. Will also set Ready for the Video struct as true.
func (p *Playlist) SetFile(vidUUID, filePath string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for s, subl := range p.playlist {
		for v, vid := range subl {
			if vid.UUID == vidUUID {
				p.playlist[s][v].File = filePath
				p.playlist[s][v].Ready = true
			}
		}
	}
}

// Returns the IP address of the Video with the provided uuid
func (p *Playlist) VideoIP(vidUUID string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, subl := range p.playlist {
		for _, vid := range subl {
			if vid.UUID == vidUUID {
				return vid.IpAddr
			}
		}
	}
	// No matching uuid
	return ""
}

// Returns the next available video that has not been played from the first
// sublist. If all the videos in the first sublist have been played, the fist
// sublist will be discarded, the remainng sublists will propagate forward
// one index, a new empty sublist will be appended to the end of the
// playlist, and the function will recurse. If the first sublist is empty,
// the function will return an empty Video struct.
func (p *Playlist) NextVideo() Video {
	p.mu.Lock()

	if len(p.playlist[0]) == 0 {
		p.mu.Unlock()
		// No Video structs in playlist, return empty Video struct.
		return Video{}
	}

	// Get unplayed Video, if available, from first sublist.
	for v, vid := range p.playlist[0] {
		if vid.Played {
			p.playlist[0][v].NP = false
		// Not played
		} else {
			if !vid.Ready {
				// Non-ready videos are still downloading
				// Wait for them to complete by returning empty
				p.mu.Unlock()
				return Video{}
			}

			p.playlist[0][v].NP = true
			p.playlist[0][v].Played = true

			p.mu.Unlock()
			return p.playlist[0][v]
		}
	}

	// Delete all videos from the completed top bucket
	for v := range p.playlist[0] {
		p.playlist[0][v].DeleteFile()
	}

	// Propagate sublists and append new empty Video sublist to end.
	p.playlist = append(p.playlist[1:], make([]Video, 0))

	p.mu.Unlock()
	// Recurse
	return p.NextVideo()
}

// Updates the Alias for each Video struct in the playlist with the given IP
// address.
func (p *Playlist) UpdateAlias(userIP, newAlias string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for s, subl := range p.playlist {
		for v, vid := range subl {
			if vid.IpAddr == userIP {
				p.playlist[s][v].Alias = newAlias
			}
		}
	}
}

// Returns a copy of the entire working playlist.
func (p *Playlist) Playlist() [][]Video {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.playlist
}
