package server

import (
	"../playlist"
)

type PageInfo struct {
	NowPlaying    playlist.Video
	Buckets       [][]playlist.Video
	ThisUserAlias string
}

func fetchInfo(addr string) PageInfo {
	var out PageInfo
	out.Buckets, out.NowPlaying = playlist.GetAllInfo()
	out.ThisUserAlias, _ = playlist.GetAlias(addr)
	return out
}
