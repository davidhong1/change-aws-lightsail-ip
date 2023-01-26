package iplist

import (
	"sync"
)

var (
	ipList []string
	rw     sync.RWMutex
)

func Get() []string {
	rw.RLock()
	defer rw.RUnlock()

	return ipList
}

func Add(ip string) {
	rw.Lock()
	defer rw.Unlock()

	for _, v := range ipList {
		if v == ip {
			return
		}
	}

	ipList = append(ipList, ip)
}

func Remove(ip string) {
	rw.Lock()
	defer rw.Unlock()

	newLen := 0
	if len(ipList) > 0 {
		newLen = len(ipList) - 1
	}
	newList := make([]string, 0, newLen)
	for _, v := range ipList {
		if v != ip {
			newList = append(newList, v)
		}
	}
	ipList = newList
}
