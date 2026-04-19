package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry holds cached secrets for a profile.
type Entry struct {
	Profile   string            `json:"profile"`
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
}

// Cache manages on-disk secret caching.
type Cache struct {
	dir string
	TTL time.Duration
}

// New returns a Cache rooted at dir with the given TTL.
func New(dir string, ttl time.Duration) *Cache {
	return &Cache{dir: dir, TTL: ttl}
}

func (c *Cache) path(profile string) string {
	return filepath.Join(c.dir, profile+".json")
}

// Get returns a cached entry if it exists and has not expired.
func (c *Cache) Get(profile string) (*Entry, bool) {
	data, err := os.ReadFile(c.path(profile))
	if err != nil {
		return nil, false
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, false
	}
	if time.Since(e.FetchedAt) > c.TTL {
		return nil, false
	}
	return &e, true
}

// Set writes an entry to disk.
func (c *Cache) Set(profile string, secrets map[string]string) error {
	if err := os.MkdirAll(c.dir, 0700); err != nil {
		return err
	}
	e := Entry{Profile: profile, Secrets: secrets, FetchedAt: time.Now()}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path(profile), data, 0600)
}

// Invalidate removes the cached entry for a profile.
func (c *Cache) Invalidate(profile string) error {
	err := os.Remove(c.path(profile))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
