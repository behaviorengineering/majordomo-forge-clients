package gitfetch

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// OriginFetchCache TTL-gates full origin fetches by common git dir.
type OriginFetchCache struct {
	mu       sync.Mutex
	now      func() time.Time
	fullAt   map[string]time.Time
	branchAt map[string]map[string]time.Time
}

// NewOriginFetchCache returns an empty fetch TTL cache.
func NewOriginFetchCache() *OriginFetchCache {
	return &OriginFetchCache{
		now:      time.Now,
		fullAt:   map[string]time.Time{},
		branchAt: map[string]map[string]time.Time{},
	}
}

// SetNow injects a clock for tests.
func (c *OriginFetchCache) SetNow(now func() time.Time) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if now == nil {
		c.now = time.Now
		return
	}
	c.now = now
}

// NeedsFullFetch reports whether a full git fetch --prune origin is required.
func (c *OriginFetchCache) NeedsFullFetch(commonDir string, ttl time.Duration, fresh bool) bool {
	commonDir = filepath.Clean(commonDir)
	if commonDir == "" {
		return true
	}
	if c == nil || ttl <= 0 || fresh {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	at, ok := c.fullAt[commonDir]
	if !ok {
		return true
	}
	return c.clock().Sub(at) >= ttl
}

func (c *OriginFetchCache) clock() time.Time {
	if c == nil || c.now == nil {
		return time.Now()
	}
	return c.now()
}

// MarkFullSuccess records a successful full prune fetch.
func (c *OriginFetchCache) MarkFullSuccess(commonDir string) {
	if c == nil {
		return
	}
	commonDir = filepath.Clean(commonDir)
	if commonDir == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fullAt == nil {
		c.fullAt = map[string]time.Time{}
	}
	c.fullAt[commonDir] = c.clock()
	c.branchAt[commonDir] = map[string]time.Time{}
}

// StaleBranches returns branches needing targeted refresh while full TTL is warm.
func (c *OriginFetchCache) StaleBranches(commonDir string, ttl time.Duration, branches []string) []string {
	commonDir = filepath.Clean(commonDir)
	cleaned := sanitizeBranchList(branches)
	if commonDir == "" || len(cleaned) == 0 {
		return nil
	}
	if c == nil || ttl <= 0 {
		return cleaned
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	byBranch := c.branchAt[commonDir]
	var stale []string
	for _, b := range cleaned {
		at, ok := byBranch[b]
		if !ok || now.Sub(at) >= ttl {
			stale = append(stale, b)
		}
	}
	return stale
}

// MarkBranchesSuccess records a successful targeted fetch.
func (c *OriginFetchCache) MarkBranchesSuccess(commonDir string, branches []string) {
	if c == nil {
		return
	}
	commonDir = filepath.Clean(commonDir)
	cleaned := sanitizeBranchList(branches)
	if commonDir == "" || len(cleaned) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.branchAt == nil {
		c.branchAt = map[string]map[string]time.Time{}
	}
	byBranch := c.branchAt[commonDir]
	if byBranch == nil {
		byBranch = map[string]time.Time{}
		c.branchAt[commonDir] = byBranch
	}
	now := c.clock()
	for _, b := range cleaned {
		byBranch[b] = now
	}
}

func sanitizeBranchList(branches []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, b := range branches {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		if _, ok := seen[b]; ok {
			continue
		}
		seen[b] = struct{}{}
		out = append(out, b)
	}
	sort.Strings(out)
	return out
}
