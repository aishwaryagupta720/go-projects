package cache

import (
	"fmt"
	"net"
	"time"
)

func InitCache() *DNSCache {

	cacheFile := "dns-cache.json"
	cache, err := NewDNSCache(cacheFile)
	if err != nil {
		fmt.Println("Error loading cache:", err)
		return nil
	}
	return cache
}

func InsertInCache(domain string, ip net.IP, ttl uint32) {

	cache := InitCache()
	// domain := "example.com"
	// ip := net.ParseIP("93.184.216.34")
	time := time.Duration(int(ttl)) * time.Second // 5 minutes

	// Add a record to the cache
	err := cache.Add(domain, ip, time)
	if err != nil {
		fmt.Println("Error adding to cache:", err)
		return
	}
}

func GetFromCache(domain string) net.IP {
	cache := InitCache()
	cachedIP := cache.Get(domain)
	if cachedIP == nil {
		// fmt.Println("No Value in Cache")
		return nil
	}

	fmt.Println("Cached IP:", cachedIP)
	return cachedIP
}
