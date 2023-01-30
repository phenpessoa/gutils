package strutils

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/bits"
	"math/rand"
	"sync"
	"unsafe"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type randomInt64 func() int64

// Rand generates a pseudo-random string with len n.
//
// The characters used to generate the string are a-zA-Z0-9.
func Rand(n int) string {
	return RandFromString(n, chars)
}

// RandFromString generates a pseudo-random string with len n.
//
// The characters used to generate the string will be retrivied from chars.
func RandFromString(n int, chars string) string {
	s := make([]byte, n)
	genRngChars(s, chars, rand.Int63)
	return *(*string)(unsafe.Pointer(&s))
}

// CryptoRand generates a cryptographically random string with len n.
//
// The characters used to generate the string are a-zA-Z0-9.
func CryptoRand(n int) string {
	return CryptoRandFromString(n, chars)
}

// CryptoRandFromString generates a cryptographically random string with len n.
//
// The characters used to generate the string will be retrivied from chars.
func CryptoRandFromString(n int, chars string) string {
	s := make([]byte, n)
	genRngChars(s, chars, cRandInt64)
	return *(*string)(unsafe.Pointer(&s))
}

var (
	cryptoRandSlice = make([]byte, 8)
	cryptoRandMu    = sync.Mutex{}
)

// cRandInt64 uses a cached slice because even though
// it is slower in parallel, it allocs way less objects
// and uses less memory.
//
// Using a cached slice is also faster in linear
// calls.
func cRandInt64() int64 {
	cryptoRandMu.Lock()
	cryptorand.Read(cryptoRandSlice)
	x := int64(binary.LittleEndian.Uint64(cryptoRandSlice))
	cryptoRandMu.Unlock()
	return x
}

// genRngChars uses the algo from this answer:
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func genRngChars(s []byte, chars string, rand randomInt64) {
	bts := bits.Len64(uint64(len(chars)))
	mask := 1<<bts - 1
	max := 63 / bts

	for i, cache, remain := len(s)-1, rand(), max; i >= 0; {
		if remain == 0 {
			cache, remain = rand(), max
		}

		if idx := int(cache & int64(mask)); idx < len(chars) {
			s[i] = chars[idx]
			i--
		}

		cache >>= bts
		remain--
	}
}
