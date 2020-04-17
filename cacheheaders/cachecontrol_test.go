package cacheheaders

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCacheControl(t *testing.T) {
	// case 1: Private content with max-age and s-maxage
	cc := &CacheControl{Private: true}
	cc.SetMaxAge(30)
	cc.SetSMaxAge(60) // s-maxage will get ignored since the content is private

	expected := "max-age=30, private"
	testCC(cc, expected, t) // tests the middleware via `router.Use(cc.SendHeaders)`

	// case 2: public, all ttl directives set
	cc = &CacheControl{}
	cc.SetMaxAge(20)
	cc.SetSMaxAge(40)
	cc.SetStaleWhileRevalidate(120)
	expected = "max-age=20, stale-while-revalidate=120, s-maxage=40"
	testCC(cc, expected, t)

	// case 3: all of the above + "no-store"
	// since NoStore is true, all other directives are ignored.
	cc = &CacheControl{NoStore: true, Private: true}
	cc.SetMaxAge(20)
	cc.SetSMaxAge(40)
	cc.SetStaleWhileRevalidate(120)
	expected = "no-store"
	testCC(cc, expected, t)
}

// testCC - tests the output of a specific cache control configuration
func testCC(cc *CacheControl, expected string, t *testing.T) {
	recorder := httptest.NewRecorder() // records output from router
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	// setup router with the CacheControl middleware and a barebones test endpoint
	router := mux.NewRouter()
	router.Use(cc.SendHeaders)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok!"))
	})

	router.ServeHTTP(recorder, request) // perform the request

	got := recorder.Header().Get("Cache-Control")

	if got != expected {
		t.Errorf("Cache Control header mismatch!\nExpected: %v\nGot     : %v\n", expected, got)
	}
}
