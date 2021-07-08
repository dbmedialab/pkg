package fasthash

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

const (
	hashLength = 24 // The string length of a base64-encoded 128-bit value
)

// TestHasher - verifies that the hasher is able to produce checksums
// Also checks the length of generated checksums vs the expected hash length.
func TestHasher(t *testing.T) {

	values := []string{
		"12",
		"Advancement",
		"I could have become a mass murderer after I hacked my governor module, but then I realized I could access the combined feed of entertainment channels carried on the company satellites.",
		"There's a reason there was such outcry when Syfy canceled The Expanse in May. The drama is one of the best sci-fi shows of recent years, and one of the best shows currently airing, period. And though the cancellation controversy and subsequent Amazon pick-up made many headlines, not enough people are actually taking the time to discuss how friggin' awesome this third season has been. Everything that has been brewing since the drama premiered is finally coming to a head in incredibly trippy fashion (wormholes! The slow zone! Ghost Thomas Jane!). In a lesser show, some of these developments might come off as cheesy or convoluted, but here they feel like a shot of juice straight to the veins as we eagerly followed the protomolecule and Holden's seemingly fated collision course.",
	}

	// 32-bit key
	key, err := randByteSlice(32)
	if err != nil {
		t.Errorf(`Key generation error: %s`, err.Error())
		return
	}

	kstr := base64.StdEncoding.EncodeToString(key)

	t.Logf("Key (b64): %#v", kstr)

	h, err := New(kstr)

	if err != nil {
		t.Errorf(`Hasher error: %s`, err.Error())
		return
	}

	for _, s := range values {

		hashStr, err := h.MakeBase64CheckSum([]byte(s))

		if err != nil {
			t.Errorf(`Checksum error: %s`, err.Error())
			return
		}
		t.Logf("HashStr: %#v", hashStr)

		l := len(hashStr)

		if l <= 16 {
			t.Errorf(`Hash string is %d bytes long, should be %d or more`, l, 16)
			return
		}

		if l > hashLength {
			t.Errorf(`Hash string is %d bytes long, should be %d or less`, l, hashLength)
			return
		}
	}
}

// TestVerifyCheckSums - Verifies that the hasher produces the expected output
func TestVerifyCheckSums(t *testing.T) {
	tests := []struct {
		key          string
		input        string
		expectedHash string
	}{
		{
			"MVyJEGNm2v5PZrCAlmblCgQAwb7F+ZzPJljAqzh+/ac=",
			"01",
			"wBXwe3qHsE9NBEkfoK5wZg==",
		},
		{
			"qHvOoDrdq4CYXGDd4UeyGG9OOfuLxdS/8F+TNrpF+xg=",
			"Ascendancy",
			"cFaBcwaL9Qjv0aLntiK0JA==",
		},
		{
			"LIa5wp1j//l4x5iZKnVMzQx5wSq65ZOla6En53zmCbU=",
			"50-53: Memory initialization error. Invalid memory type or incompatible memory speed.",
			"7vC9U5tpLn4OvzoXybBzMw==",
		},
		{
			"y6pghJ0clnqqeACueXC+KsFwVQ1X6k4tK6he0T9I0IY=",
			"Magic is unknown, dangerous and inhuman. Even the best wizards occasionally fail to properly harness a spell, with unpredictable results. Wizards thus inculcate their preferred magics, lest they err in casting a spell and corrupt themselves with misdirected magical energies. At 1st level a wizard determines 4 spells that he knows, representing years of study and practice. As his comprehension expands, a wizard may learn more spells of progressively higher levels. A wizard knows a number of spells as shown on table 1-12, modified by his Intelligence score.",
			"EbnR3vHy0395DygyiTv0sw==",
		},
	}

	for _, ct := range tests {
		h, err := New(ct.key)
		if err != nil {
			t.Fatalf("Failed to create hasher: %s", err.Error())
		}

		hashStr, err := h.MakeBase64CheckSum([]byte(ct.input))
		if err != nil {
			t.Fatalf("Failed to generate hash: %s", err.Error())
		}

		if hashStr != ct.expectedHash {
			t.Errorf("Hash mismatch!\nExpected: %s\nGot:      %s\n", ct.expectedHash, hashStr)
		}
	}
}

// TestThreadSafety - Checks if it's safe to simultaneously use an instance of the hasher from multiple threads
// Reusing an instance won't produce huge performance gains, but it will save you a few lines of code.
func TestThreadSafety(t *testing.T) {
	tests := []threadTest{}

	key, err := randByteSlice(32)
	if err != nil {
		t.Errorf(`Key generation error: %s`, err.Error())
		return
	}

	kstr := base64.StdEncoding.EncodeToString(key)

	// initialize hasher
	h, err := New(kstr)
	if err != nil {
		t.Errorf("failed to initialize hasher: %s", err.Error())
	}

	// generate 10000 tests
	for i := 0; i < 10000; i++ {
		tt := threadTest{}
		bs, err := randByteSlice(64)
		if err != nil {
			t.Fatalf(`Input generation error: %s`, err.Error())
		}

		tt.input = bs
		hashStr, err := h.MakeBase64CheckSum(bs)
		if err != nil {
			t.Fatalf("Failed to generate hash: %s", err.Error())
		}

		tt.expectedHash = hashStr
		tests = append(tests, tt)
	}

	// Try to replicate the results, using the same hasher instance simultaneously from 10000 threads.
	for _, tt := range tests {
		go tt.run(h, t)
	}
}

type threadTest struct {
	input        []byte
	expectedHash string
}

func (tt threadTest) run(h *Hasher, t *testing.T) {
	hashStr, err := h.MakeBase64CheckSum(tt.input)
	if err != nil {
		t.Errorf("Failed to generate hash: %s", err.Error())
	}

	if hashStr != tt.expectedHash {
		t.Errorf("Hash mismatch!\nExpected: %s\nGot:      %s\n", tt.expectedHash, hashStr)
	}
}

// randByteSlice - generates a random byteslice of length `len`
func randByteSlice(len int) (b []byte, err error) {
	b = make([]byte, len)
	_, err = rand.Read(b)
	return
}
