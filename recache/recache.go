// Package recache provides a thread-safe LRU cache of compiled regular
// expressions.
package recache

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/charlievieth/reonce"
)

type entry struct {
	next, prev *entry
	re         *reonce.Regexp
}

// Cache is a LRU cache of compiled Regexps. All methods are safe for
// concurrent access. The zero value for Cache is an empty non-POSIX
// cache with no max size and if safe for use.
type Cache struct {
	mu         sync.Mutex
	cache      map[string]*entry
	ll         *list
	maxEntries int // zero means no limit
	posix      bool
}

func newCache(maxEntries int, posix bool) *Cache {
	if maxEntries < 0 {
		panic("recache: non-positive maxEntries: " + strconv.Itoa(maxEntries))
	}
	return &Cache{maxEntries: maxEntries, posix: posix}
}

// New creates a new LRU Cache that will cache maxEntries Regexps.
// If maxEntries is zero there is no limit. New panics if maxEntries
// if less than zero.
func New(maxEntries int) *Cache { return newCache(maxEntries, false) }

// New creates a new LRU Cache that for POSIX Regexps. The Compile and
// MustCompile methods call regexp.Compile and regexp.MustCompile.
func NewPOSIX(maxEntries int) *Cache { return newCache(maxEntries, true) }

func (c *Cache) POSIX() bool { return c.posix }

// MaxEntries returns the maximum size of the Cache.
func (c *Cache) MaxEntries() int {
	c.mu.Lock()
	n := c.maxEntries
	c.mu.Unlock()
	return n
}

// SetMaxEntries sets the maximum number of Regexps that will be cached before
// an item is evicted and return the previous max. If n is smaller than the
// current number of cached entries, the cache is trimmed. If n == 0 the max
// number of entries is infinite. If n < 0 SetMaxEntries panics.
func (c *Cache) SetMaxEntries(n int) (prev int) {
	if n < 0 {
		panic("recache: non-positive value n: " + strconv.Itoa(n))
	}
	c.mu.Lock()
	if n != 0 && c.ll != nil {
		for i := c.ll.Len() - n; i > 0; i-- {
			c.removeOldest()
		}
	}
	prev = c.maxEntries
	c.maxEntries = n
	c.mu.Unlock()
	return prev
}

func (c *Cache) lazyInit() {
	c.cache = make(map[string]*entry)
	c.ll = newList()
}

func (c *Cache) get(expr string) *reonce.Regexp {
	c.mu.Lock()
	ee := c.cache[expr]
	if ee != nil {
		c.ll.MoveToFront(ee)
	} else {
		if c.cache == nil {
			c.lazyInit()
		}
		if !c.posix {
			ee = &entry{re: reonce.New(expr)}
		} else {
			ee = &entry{re: reonce.NewPOSIX(expr)}
		}
		if c.maxEntries != 0 && c.ll.Len() >= c.maxEntries {
			c.removeOldest()
		}
		c.cache[expr] = c.ll.PushFront(ee)
	}
	c.mu.Unlock()
	return ee.re
}

// Compile compiles the Regexp and panics if there is an error.
// If the Regexp has already been compiled the cached Regexp is returned.
// Otherwise the Regexp is compiled and added to the Cache.
func (c *Cache) Compile(key string) (*regexp.Regexp, error) {
	re := c.get(key)
	if err := re.Compile(); err != nil {
		return nil, err
	}
	return re.Regexp(), nil
}

// MustCompile compiles the Regexp and panics if there is an error.
// If the Regexp has already been compiled the cached Regexp is returned.
func (c *Cache) MustCompile(key string) *regexp.Regexp {
	return c.get(key).Regexp()
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
	delete(c.cache, e.re.String())
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	n := len(c.cache)
	c.mu.Unlock()
	return n
}

// DefaultCacheSize is the size of default Cache.
const DefaultCacheSize = 256

var (
	std   = New(DefaultCacheSize)
	posix = NewPOSIX(DefaultCacheSize)
)

// Compile compiles the Regexp and panics if there is an error.
// If the Regexp has already been compiled the cached Regexp is returned.
// Otherwise the Regexp is compiled and added to the default cache.
func Compile(str string) (*regexp.Regexp, error) {
	return std.Compile(str)
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular
// expressions.
func MustCompile(str string) *regexp.Regexp {
	return std.MustCompile(str)
}

// CompilePOSIX compiles the Regexp and panics if there is an error.
// If the Regexp has already been compiled the cached Regexp is returned.
// Otherwise the Regexp is compiled and added to the default cache.
func CompilePOSIX(str string) (*regexp.Regexp, error) {
	return posix.Compile(str)
}

// MustCompilePOSIX is like CompilePOSIX but panics if the expression cannot
// be parsed. It simplifies safe initialization of global variables holding
// compiled regular expressions.
func MustCompilePOSIX(str string) *regexp.Regexp {
	return posix.MustCompile(str)
}

// MaxEntries returns the size of the default Cache.
func MaxEntries() int {
	return std.MaxEntries()
}

// SetMaxEntries changes the size of the default Cache.
func SetMaxEntries(n int) (prev int) {
	return std.SetMaxEntries(n)
}

// Len returns the number of cached Regexps in the default Cache.
func Len() int { return std.Len() }

// MaxEntriesPOSIX returns the size of the default POSIX Cache.
func MaxEntriesPOSIX() int {
	return posix.MaxEntries()
}

// SetMaxEntriesPOSIX changes the size of the default POSIX Cache.
func SetMaxEntriesPOSIX(n int) (prev int) {
	return posix.SetMaxEntries(n)
}

// LenPOSIX returns the number of cached Regexps in the default POSIX Cache.
func LenPOSIX() int { return posix.Len() }
