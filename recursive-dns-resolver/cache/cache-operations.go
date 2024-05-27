package cache

import (
	"encoding/json"
	"net"
	"os"
	"time"
)

func NewDNSCache(file string) (*DNSCache, error) {
	cache := &DNSCache{
		Records: make(map[string]CacheRecord),
		file:    file,
	}

	err := cache.loadFromFile()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// Load records from the file
func (cache *DNSCache) loadFromFile() error {
	data, err := os.ReadFile(cache.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, start with an empty cache
		}
		return err
	}

	return json.Unmarshal(data, &cache.Records)
}

// Save records to the file
func (cache *DNSCache) saveToFile() error {
	data, err := json.MarshalIndent(cache.Records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cache.file, data, 0644)
}

// Add a record to the cache
func (cache *DNSCache) Add(domain string, ip net.IP, ttl time.Duration) error {
	cache.Records[domain] = CacheRecord{
		Domain:    domain,
		IPAddress: ip,
		ExpiresAt: time.Now().Add(ttl),
	}

	return cache.saveToFile()
}

// Get a record from the cache, removing expired records
func (cache *DNSCache) Get(domain string) net.IP {
	cache.removeExpired()
	record, exists := cache.Records[domain]
	if !exists {
		return nil
	}
	return record.IPAddress
}

// Remove expired records from the cache
func (cache *DNSCache) removeExpired() {

	now := time.Now()
	for domain, record := range cache.Records {
		if record.ExpiresAt.Before(now) {
			delete(cache.Records, domain)
		}
	}
	cache.saveToFile() // Save the updated cache to the file
}
