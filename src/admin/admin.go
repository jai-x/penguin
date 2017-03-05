package admin

import (
	"fmt"
	"time"
	"sync"
	"crypto/sha256"

	"../help"
	"../config"
)

type AdminSessions struct {
	// Map of admin ip addresses to session timeout token
	Admins map[string]time.Time
	AdminLock sync.RWMutex

	// SHA256 hash of pwd
	pwdHash string
}

// Constructor
func (a *AdminSessions) Init() {
	// Init session timeout map
	a.Admins = make(map[string]time.Time)

	// Get password, hash and store
	h := sha256.New()
	h.Write([]byte(config.Config.AdminPass))
	a.pwdHash = fmt.Sprintf("%x", h.Sum(nil))
}

// Get session from client address
func (a *AdminSessions) ValidSession(addr string) bool {
	ip := help.GetIP(addr)

	a.AdminLock.RLock()
	defer a.AdminLock.RUnlock()

	expiry, exists := a.Admins[ip]
	if !exists || time.Now().After(expiry) {
		// Invalid session
		return false
	} else {
		return true
	}
}

// start session given client address
func (a *AdminSessions) StartSession(addr string) {
	ip := help.GetIP(addr)

	a.AdminLock.Lock()
	defer a.AdminLock.Unlock()

	// Timeout is 1 hour later
	a.Admins[ip] = time.Now().Add(1 * time.Hour)
}

// End session at a given address
func (a *AdminSessions) EndSession(addr string) {
	ip := help.GetIP(addr)

	a.AdminLock.Lock()
	defer a.AdminLock.Unlock()

	// Delete session from map
	delete(a.Admins, ip)
}

// Verify given password against password hash
func (a *AdminSessions) VerifyPassword(pass string) bool {
	h := sha256.New()
	h.Write([]byte(pass))
	return a.pwdHash == fmt.Sprintf("%x", h.Sum(nil))
}