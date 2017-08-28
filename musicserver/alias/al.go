package alias

import (
	"sync"
)

type AliasMgr struct {
	mu      sync.RWMutex
	aliases map[string]string
}

func NewAliasMgr() *AliasMgr {
	out := AliasMgr{}

	out.aliases = make(map[string]string)
	return &out
}

func (a *AliasMgr) Alias(ip string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	x, y := a.aliases[ip]
	return x, y
}

func (a *AliasMgr) SetAlias(ip, newAlias string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.aliases[ip] = newAlias
}
