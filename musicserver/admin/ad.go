package admin

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type AdminSessions struct {
	mu       sync.RWMutex
	sessions map[string]time.Time

	pwdHash string
}

func NewAdminSessions(newPwd string, preHash bool) AdminSessions {
	out := AdminSessions{}

	out.sessions = make(map[string]time.Time)

	// newPwd may be a hashed value already
	if preHash {
		out.pwdHash = newPwd
	} else {
		out.pwdHash = hashString(newPwd)
	}

	return out
}

func (a *AdminSessions) ValidSession(ip string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	expiry, exists := a.sessions[ip]

	if !exists || time.Now().After(expiry) {
		// Invalid session
		return false
	} else {
		return true
	}
}

func (a *AdminSessions) StartSession(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.sessions[ip] = time.Now().Add(1 * time.Hour)
}

func (a *AdminSessions) EndSession(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Delete session from map
	delete(a.sessions, ip)
}

func (a AdminSessions) ValidPassword(pwd string) bool {
	return a.pwdHash == hashString(pwd)
}

func hashString(in string) string {
	hasher := sha256.New()
	hasher.Write([]byte(in))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
