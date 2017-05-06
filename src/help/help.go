package help

import (
	"log"
	"fmt"
	"strings"
	"crypto/rand"
)

// Filled with commonly used convinience functions

// Addresses are in form "xxx.xxx.xxx.xxx:port"
// This strips the port number, returning only the IP
func GetIP(addr string) string {
	in := strings.LastIndex(addr, ":")
	return addr[:in]
}

// Generate a pseudo random guid
// this is the only thing Go doesnt have in a standard lib
// TODO: Find a better way to do this
func GenUUID() string {
	// Fill slice b with random bytes
	b := make([]byte, 16)
	_, err := rand.Read(b)
	// Error check
	if err != nil {
		log.Fatal("Error: ", err)
	}
	// Print to variable
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func StripYoutubePlaylist(link string) string {
	if strings.Contains(link, "youtube.com") && strings.Contains(link,"&list") {
		return strings.Split(link, "&list")[0]
	} else {
		return link
	}
}

func PrintMasthead() {
	fmt.Println("                              _            ")
	fmt.Println(" _ __   ___ _ __   __ _ _   _(_)_ __      ")
	fmt.Println("| '_ \\ / _ \\ '_ \\ / _` | | | | | '_ \\ ")
	fmt.Println("| |_) |  __/ | | | (_| | |_| | | | | |    ")
	fmt.Println("| .__/ \\___|_| |_|\\__, |\\__,_|_|_| |_| ")
	fmt.Println("|_|               |___/                   \n")
}
