package admin

import (
	"fmt"
	"time"
	"sync"
	"crypto/sha256"

	"../help"
)

func Expired(sess time.Time) bool {
	return time.Now().After(sess)
}

type AdminSessions struct {
	// Map of admin ip addresses to session timeout token
	Admins map[string]time.Time
	AdminLock sync.RWMutex

	// SHA256 hash of pwd
	pwdHash string
}

// Constructor
func (a *AdminSessions) Init(pass string) {
	a.Admins = make(map[string]time.Time)

	h := sha256.New()
	h.Write([]byte(pass))
	a.pwdHash = fmt.Sprintf("%x", h.Sum(nil))
}

// Get session from client address
func (a *AdminSessions) GetSession(addr string) (time.Time, bool){
	ip := help.GetIP(addr)

	a.AdminLock.RLock()
	defer a.AdminLock.RUnlock()

	expiry, exists := a.Admins[ip]
	return expiry, exists
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