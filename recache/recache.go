package recache

import (
	"regexp"
	"strconv"
	"sync"
)

type entry struct {
	next, prev *entry
	once       sync.Once
	posix      bool
	re         *regexp.Regexp
	err        error
	expr       string
}

var (
	compile      = regexp.Compile
	compilePOSIX = regexp.CompilePOSIX
)

func (e *entry) Compile() (*regexp.Regexp, error) {
	e.once.Do(func() {
		if e.posix {
			e.re, e.err = compilePOSIX(e.expr)
		} else {
			e.re, e.err = compile(e.expr)
		}
	})
	return e.re, e.err
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func (e *entry) MustCompile() *regexp.Regexp {
	re, err := e.Compile()
	if err != nil {
		if e.posix {
			panic(`regexp: CompilePOSIX(` + quote(e.expr) + `): ` + e.err.Error())
		} else {
			panic(`regexp: Compile(` + quote(e.expr) + `): ` + e.err.Error())
		}
	}
	return re
}

type cacheKey struct {
	expr  string
	posix bool
}

// Cache is a LRU cache of compiled regexes. It is safe for concurrent access.
type Cache struct {
	mu         sync.Mutex
	cache      map[cacheKey]*entry
	ll         *list
	maxEntries int // zero means no limit
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New(maxEntries int) *Cache {
	return &Cache{
		maxEntries: maxEntries,
		ll:         newList(),
		cache:      make(map[cacheKey]*entry),
	}
}

func (c *Cache) MaxEntries() int {
	c.mu.Lock()
	n := c.maxEntries
	c.mu.Unlock()
	return n
}

func (c *Cache) SetMaxEntries(n int) (prev int) {
	c.mu.Lock()
	if n > 0 {
		// TODO: test this
		for i := c.ll.Len() - n; i > 0; i-- {
			c.removeOldest()
		}
	}
	prev = c.maxEntries
	c.maxEntries = n
	c.mu.Unlock()
	return prev
}

func (c *Cache) get(expr string, posix bool) *entry {
	key := cacheKey{expr, posix}
	c.mu.Lock()
	ee := c.cache[key]
	if ee != nil {
		c.ll.MoveToFront(ee)
	} else {
		if c.maxEntries != 0 && c.ll.Len() >= c.maxEntries {
			c.removeOldest()
		}
		ee = c.ll.PushFront(&entry{expr: expr, posix: posix})
		c.cache[key] = ee
	}
	c.mu.Unlock()
	return ee
}

func (c *Cache) Compile(key string) (*regexp.Regexp, error) {
	return c.get(key, false).Compile()
}

func (c *Cache) MustCompile(key string) *regexp.Regexp {
	return c.get(key, false).MustCompile()
}

func (c *Cache) CompilePOSIX(key string) (*regexp.Regexp, error) {
	return c.get(key, true).Compile()
}

func (c *Cache) MustCompilePOSIX(key string) *regexp.Regexp {
	return c.get(key, true).MustCompile()
}

// removeOldest removes the oldest item from the cache.
func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *entry) {
	c.ll.Remove(e)
	kv := cacheKey{e.expr, e.posix}
	delete(c.cache, kv)
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	n := len(c.cache)
	c.mu.Unlock()
	return n
}

const DefaultCacheSize = 256

var std = New(DefaultCacheSize)

func Compile(str string) (*regexp.Regexp, error) {
	return std.Compile(str)
}

func MustCompile(str string) *regexp.Regexp {
	return std.MustCompile(str)
}

func MaxEntries() int {
	return std.MaxEntries()
}

func SetMaxEntries(n int) (prev int) {
	return std.SetMaxEntries(n)
}

func CompilePOSIX(str string) (*regexp.Regexp, error) {
	return std.CompilePOSIX(str)
}

func MustCompilePOSIX(str string) *regexp.Regexp {
	return std.MustCompilePOSIX(str)
}
