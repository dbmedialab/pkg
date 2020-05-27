# Fasthash

An easy and fast way to generate validation checksums for non-sensitive data.

It's basically a thin abstraction layer on top of the `minio/highwayhash` hash implementation, with some unit tests to verify that it works appropriately.

If you need a key, run `go test -v .` in this package. The first log line will contain a fresh key generated via "crypto/rand".

**NB:** Don't use this to hash passwords. Don't.

#### Usage Example
```go
package main

import (
	"log"

	"github.com/dbmedialab/pkg/fasthash"
)

const (
	key   = "KZyaV28e67wtgRMN6QQtX0wUhB9lj8qYDCISpOmxgKY="
	input = "Build spacecraft, fly them, and try to help the Kerbals fulfill their ultimate mission of conquering space."
)

func main() {
	h, err := fasthash.New(key)
	if err != nil {
		log.Printf("Error while initializing hasher: %s", err.Error())
		return
	}

	hashStr, err := h.MakeBase64CheckSum([]byte(input))
	if err != nil {
		log.Printf(`Checksum error: %s`, err.Error())
		return
	}

	log.Printf("Hash: %s", hashStr)
}
```