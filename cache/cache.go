package cache

import (
	"regexp"

	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/thoas/go-funk"
)

type RegexCache struct {
	simplelru.LRUCache
}

func NewRegexCache(lru simplelru.LRUCache) *RegexCache {
	return &RegexCache{
		LRUCache: lru,
	}
}

// Adds a regular expression key to the cache, returns true if an eviction
// occurred and updates the "recently used"-ness of the key.
func (r *RegexCache) DefineSysRegexp(key *regexp.Regexp, value interface{}) bool {
	return r.LRUCache.Add(key.String(), value)
}

// Updates a regular expression key and maintains the data associated
// to the old key. #updated
func (r *RegexCache) UpdateSysRegexp(key, newKey *regexp.Regexp) bool {
	if value, isFound := r.LRUCache.Get(key.String()); isFound {
		r.LRUCache.Add(newKey.String(), value)
		r.LRUCache.Remove(key.String())
		return true
	}
	return false
}

// Adds a value to the cache, returns true if an eviction occurred and
// updates the "recently used"-ness of the key.
func (r *RegexCache) Add(key string, value interface{}) bool {
	if re, found := r.matchOldest(key); found {
		return r.LRUCache.Add(re.String(), value)
	}
	return false
}

// Returns key's value from the cache and
// updates the "recently used"-ness of the key. #value, isFound
func (r *RegexCache) Get(key string) (interface{}, bool) {
	if re, found := r.matchNewest(key); found {
		return r.LRUCache.Get(re.String())
	}
	return nil, false
}

// Checks if a key exists in cache without updating the recent-ness.
func (r *RegexCache) Contains(key string) bool {
	_, found := r.matchOldest(key)
	return found
}

// Returns key's value without updating the "recently used"-ness of the key.
func (r *RegexCache) Peek(key string) (interface{}, bool) {
	if re, found := r.matchNewest(key); found {
		return r.LRUCache.Peek(re.String())
	}
	return nil, false
}

func (r *RegexCache) matchOldest(key string) (*regexp.Regexp, bool) {
	return matchRegExp(key, r.Keys())
}

func (r *RegexCache) matchNewest(key string) (*regexp.Regexp, bool) {
	return matchRegExp(key, funk.Reverse(r.Keys()).([]interface{}))
}

func matchRegExp(key string, patterns []interface{}) (*regexp.Regexp, bool) {
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.(string))
		if re.MatchString(key) {
			return re, true
		}
	}
	return nil, false
}
