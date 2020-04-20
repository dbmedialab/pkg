# Cache headers

This package can be used to set up HTTP middleware that handles the
`Cache-Control`, `X-Cache-Channel` and `Cache-Tag` headers for you.  
It is part utility and part code-as-documentation,
with the end goal of making it very difficult to mess up your cache headers.

The middleware is compatible with `gorilla/mux` and `go-chi/chi`, as well as
any other router that uses the `net/http.HandlerFunc` middleware standard.

### Usage example
```go
package server

import(
    "github.com/dbmedialab/pkg/cacheheaders"
    "github.com/gorilla/mux"
)

// addCacheHeaders - Adds cache middleware to a gorilla/mux router,
// which sends a ttl response (as cache-control directives) and a couple of cache tag / channel headers
func addCacheHeaders(r *mux.Router, ttl int, channels ...string) {
    ctrl := &cacheheaders.CacheControl{}
    ctrl.SetMaxAge(ttl)
    ctrl.SetSMaxAge(ttl)
    r.Use(ctrl.SendHeaders)
    
    chans := &cacheheaders.CacheChannels{Varnish: true, Cloudflare: true}
    chans.Set(channels...)
    r.Use(chans.SendHeaders)
}

// SetupRoutes - will get called from the main cmd to setup this server's routes
func SetupRoutes(r *mux.Router) {
    // This specific song has its own route for some reason,
    // as well as a lower cache TTL and its own cache channel
    specific := r.PathPrefix("/fischerspooner/emerge").Subrouter()
    addCacheHeaders(pluss, 30, "fischerspooner", "emerge")

    // The rest of the discography gets no special treatment
    pluss := r.PathPrefix("/fischerspooner").Subrouter()
    addCacheHeaders(pluss, 300, "fischerspooner")
}

```
