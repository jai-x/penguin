package admin

import (
	"fmt"
	"log"
	"time"
	"sync"
	"crypto/sha256"

	"../help"
	"../config"
)

var (
	// Map of admin ip addresses to session timeout token
	adminLock sync.RWMutex
	adminMap map[string]time.Time

	// SHA256 hash of pwd
	pwdHash string
)

// Constructor
func Init() {
	log.Println("Admin init...")
	// Init session timeout map
	adminMap = make(map[string]time.Time)

	// Get password, hash and store
	h := sha256.New()
	h.Write([]byte(config.Config.AdminPass))
	pwdHash = fmt.Sprintf("%x", h.Sum(nil))
}

// Get session from client address
func ValidSession(addr string) bool {
	ip := help.GetIP(addr)

	adminLock.RLock()
	defer adminLock.RUnlock()

	expiry, exists := adminMap[ip]
	if !exists || time.Now().After(expiry) {
		// Invalid session
		return false
	} else {
		return true
	}
}

// start session given client address
func StartSession(addr string) {
	ip := help.GetIP(addr)

	adminLock.Lock()
	defer adminLock.Unlock()

	// Timeout is 1 hour later
	adminMap[ip] = time.Now().Add(1 * time.Hour)
}

// End session at a given address
func EndSession(addr string) {
	ip := help.GetIP(addr)

	adminLock.Lock()
	defer adminLock.Unlock()

	// Delete session from map
	delete(adminMap, ip)
}

// Verify given password against password hash
func VerifyPassword(pass string) bool {
	h := sha256.New()
	h.Write([]byte(pass))
	return pwdHash == fmt.Sprintf("%x", h.Sum(nil))
}
