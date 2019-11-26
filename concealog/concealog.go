// Package concealog : can be used to remove credentials from request log data
package concealog

import (
	"regexp"
)

// AuthReplacer - Replaces the `Authorization` header in request dumps
type AuthReplacer struct {
	cReg *regexp.Regexp
	cRep string
}

// NewAuthReplacer - returns an initialized AuthReplacer
func NewAuthReplacer() (ar *AuthReplacer, err error) {
	ar = &AuthReplacer{}
	ar.cReg, err = regexp.Compile(`Authorization: (\w+) [\w-]+`)
	ar.cRep = `Authorization: $1 *********`
	return
}

// ReplaceString - takes a request body as input,
// and replaces any authorization headers with stars
func (ar *AuthReplacer) ReplaceString(body string) string {
	return ar.cReg.ReplaceAllString(body, ar.cRep)
}
