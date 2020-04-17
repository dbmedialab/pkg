package cacheheaders

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCacheChannels(t *testing.T) {
	// Test 1: Strip non-alphanumeric characters
	cc := &CacheChannels{
		Varnish: true,
	}
	cc.Set("Cat(øøøøø)_Articles//")

	expected := cacheChanExpected{Channel: "Cat_Articles"}

	testCChan(cc, expected, t)

	// Test 2: Comma-separation (without spaces on cloudflare)
	cc.Add("Cat_Pictures", "Dog_Pictures")
	cc.Cloudflare = true
	expected = cacheChanExpected{
		Channel: "Cat_Articles, Cat_Pictures, Dog_Pictures",
		Tag:     "Cat_Articles,Cat_Pictures,Dog_Pictures",
	}
}

type cacheChanExpected struct {
	Channel string
	Tag     string
}

// testCChan - tests the output of a specific cache channel configuration
func testCChan(cc *CacheChannels, expected cacheChanExpected, t *testing.T) {
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

	gotChan := recorder.Header().Get(CacheChannelHeader)
	if gotChan != expected.Channel {
		t.Errorf("Cache Channel header mismatch!\nExpected: %v\nGot     : %v\n", expected.Channel, gotChan)
	}

	gotTag := recorder.Header().Get(CacheTagHeader)
	if gotTag != expected.Tag {
		t.Errorf("Cache Tag header mismatch!\nExpected: %v\nGot     : %v\n", expected.Tag, gotTag)
	}
}
