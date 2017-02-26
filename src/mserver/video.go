package mserver

type Video struct {
	Title string
	File string
	IpAddr string
}

type VideoInfo struct {
	Title string
	Uploader string
}

func (v Video) ConvertToInfo(aliasMap map[string]string) VideoInfo {
	name, exists := aliasMap[v.IpAddr]
	if exists {
		return VideoInfo{v.Title, name}
	} else {
		return VideoInfo{v.Title, "Anon"}
	}
}