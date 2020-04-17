package cacheheaders

import (
	"net/http"
	"strconv"
	"strings"
)

// CacheControl - A middleware struct that outputs cache control directives for browsers and cache proxies
// Currently limited to the core set of directives necessary to achieve configurable caching and TTL.
// NB: This implementation assumes that the CacheControl configuration will not change during runtime.
type CacheControl struct {
	Private              bool    // "private"  - tells Varnish and Cloudflare to never cache this response (browsers will though)
	NoStore              bool    // "no-store" - tells proxy caches and browsers that this response contains sensitive data and should not be cached or stored anywhere after use
	maxAge               *int    // see SetMaxAge
	sMaxAge              *int    // see SetSMaxAge
	staleWhileRevalidate *int    // see SetStaleWhileRevalidate
	cached               *string // contains the cache-control string after first generation
}

// SetMaxAge - Sets the "max-age" header:
// The TTL (in seconds) that browsers should obey. Varnish will use this if s-maxage is not set.
func (cc *CacheControl) SetMaxAge(value int) {
	cc.maxAge = &value
}

// SetSMaxAge - Sets the "s-maxage" header:
// The TTL (in seconds) that Varnish, Cloudflare and other proxy-caches should obey
func (cc *CacheControl) SetSMaxAge(value int) {
	cc.sMaxAge = &value
}

// SetStaleWhileRevalidate - Sets the "stale-while-revalidate" header:
// The number of seconds during which browsers will reuse a stale response while sending a revalidation request in the background
func (cc *CacheControl) SetStaleWhileRevalidate(value int) {
	cc.staleWhileRevalidate = &value
}

// SendHeaders - A middleware function compatible with most routers.
// Outputs cache headers according to the CacheControl configuration.
func (cc *CacheControl) SendHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ccs := cc.cacheControlString()
		if ccs != "" {
			w.Header().Set("Cache-Control", ccs)
		}

		next.ServeHTTP(w, r)
	})
}

func (cc *CacheControl) cacheControlString() string {
	if cc.cached != nil {
		return *cc.cached // string has already been generated once
	}

	// make a slice of cache control directives
	ccSlice := cc.makeSlice()

	// join directives into a string, then cache it
	joined := strings.Join(ccSlice, ", ")
	cc.cached = &joined

	return joined
}

// makeSlice - constructs a slice of cache-control headers based on CacheControl config
func (cc *CacheControl) makeSlice() []string {

	if cc.NoStore {
		return []string{"no-store"}
	}

	cacheControl := []string{}

	if cc.maxAge != nil {
		cacheControl = append(cacheControl, "max-age="+strconv.Itoa(*cc.maxAge))
	}

	if cc.staleWhileRevalidate != nil {
		cacheControl = append(cacheControl, "stale-while-revalidate="+strconv.Itoa(*cc.staleWhileRevalidate))
	}

	if cc.Private {
		cacheControl = append(cacheControl, "private")
		return cacheControl // proxy caching disallowed: No need to send proxy-specific headers
	}

	if cc.sMaxAge != nil {
		cacheControl = append(cacheControl, "s-maxage="+strconv.Itoa(*cc.sMaxAge))
	}

	return cacheControl
}
