package cache

import (
	"net"
	"sync"
	"time"
)

type CacheRecord struct {
	Domain    string
	IPAddress net.IP
	ExpiresAt time.Time
}

type DNSCache struct {
	Records map[string]CacheRecord
	mutex   sync.Mutex
	file    string
}
