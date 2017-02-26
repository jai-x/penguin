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
