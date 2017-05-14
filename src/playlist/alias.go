package playlist

import (
	"sync"

	"../help"
)

var (
	aliasLock sync.RWMutex
	aliasMap  map[string]string
)

func GetAlias(addr string) (string, bool) {
	aliasLock.RLock()
	defer aliasLock.RUnlock()

	ip := help.GetIP(addr)
	val, ok := aliasMap[ip]
	return val, ok
}

func SetAlias(addr, alias string) {
	aliasLock.Lock()
	defer aliasLock.Unlock()

	ip := help.GetIP(addr)
	aliasMap[ip] = alias

	go updateAliases()
}

func updateAliases() {
	aliasLock.RLock()
	defer aliasLock.RUnlock()

	bucketLock.Lock()
	defer bucketLock.Unlock()

	/* Update all alias values from the ip of each video struct against the
	 * alias map */
	for b, bucket := range buckets {
		for v, vid := range bucket {
			buckets[b][v].Alias, _ = aliasMap[vid.IpAddr]
		}
	}

	stateChange()
}
