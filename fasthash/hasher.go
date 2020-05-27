// package fasthash - A utility that provides a quick hash function,
// handy for generating validation checksums for non-private data.
// It's basically a thin abstraction layer on top of the `minio/highwayhash` hash implementation,
// with some unit tests to verify that it works appropriately.
//
// If you need a key, run `go test -v .` in this package.
// The first log line will contain a fresh key generated via "crypto/rand".
//
// NB: Don't use this to hash passwords. Don't.
package fasthash

import (
	"encoding/base64"
	"github.com/minio/highwayhash"
)

// Hasher - The main hasher struct.
// Can be used by multiple threads simultaneously (as long as none of them change the key).
type Hasher struct {
	base64Key string // base64-encoded 32-byte hash key
	byteKey   []byte
}

// New - Takes a base64-encoded 32-byte key string, and returns an initialized hasher.
func New(key string) (h *Hasher, err error) {
	h = &Hasher{}
	err = h.applyKey(key)
	return
}

// applyKey - Decodes the provided base64-encoded 32-byte key and applies it to the hasher.
func (h *Hasher) applyKey(key string) error {
	h.base64Key = key
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	h.byteKey = k
	return nil
}

// MakeBase64CheckSum - Returns a base64-encoded checksum based on the input byteslice.
func (h *Hasher) MakeBase64CheckSum(b []byte) (s string, err error) {
	hash, err := highwayhash.New128(h.byteKey)
	if err != nil {
		return
	}
	hash.Write(b)
	sum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum), nil
}
