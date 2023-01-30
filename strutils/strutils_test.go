package strutils

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"testing"
	"unsafe"
)

const testStrLen = 100

var sink string

func cryptoRandNoCache(n int) string {
	return cryptoRandFromStringNoCache(n, chars)
}

func cryptoRandFromStringNoCache(n int, chars string) string {
	s := make([]byte, n)
	genRngChars(s, chars, cRandInt64NoCache)
	return *(*string)(unsafe.Pointer(&s))
}

func cRandInt64NoCache() int64 {
	cryptoRandSlice := make([]byte, 8)
	cryptorand.Read(cryptoRandSlice)
	return int64(binary.LittleEndian.Uint64(cryptoRandSlice))
}

func Benchmark_CryptoRand(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = CryptoRand(testStrLen)
	}
}

func Benchmark_CryptoRand_NoCache(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = cryptoRandNoCache(testStrLen)
	}
}

func Benchmark_CryptoRand_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			sink = CryptoRand(testStrLen)
		}
	})
}

func Benchmark_CryptoRand_Parallel_NoCache(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			sink = cryptoRandNoCache(testStrLen)
		}
	})
}

func testDistribution(f func(int) string) error {
	const (
		distLoopN = 10000
		accept    = distLoopN / 100 // 1% max difference for each character
		n         = 1
	)

	counter := make(map[byte]int, len(chars))

	for i := 0; i < distLoopN; i++ {
		str := f(n)
		counter[str[0]]++
	}

	max := -1
	min := math.MaxInt
	for _, c := range counter {
		if c < min {
			min = c
		}
		if c > max {
			max = c
		}
	}

	if (max - min) > accept {
		return fmt.Errorf(
			"\nbiased result:\nmap: %#v\nmax: %d\nmin: %d\ndiff: %d\naccept: %d",
			counter, max, min, max-min, accept,
		)
	}

	return nil
}

func TestPseudoDistribution(t *testing.T) {
	if err := testDistribution(Rand); err != nil {
		t.Error(err)
	}
}

func TestCryptoDistribution(t *testing.T) {
	if err := testDistribution(CryptoRand); err != nil {
		t.Error(err)
	}
}

func testRepetition(f func(int) string) error {
	const (
		strLen = 16
		tries  = 10000
		accept = 0
	)

	type token struct{}
	cache := make(map[string]token)

	for i := 0; i < tries; i++ {
		cache[f(strLen)] = token{}
	}

	if (tries - len(cache)) > accept {
		return fmt.Errorf(
			"\ntoo many repetitions:\naccept: %d\ngot: %d",
			accept, tries-len(cache),
		)
	}

	return nil
}

func TestPseudoRepetition(t *testing.T) {
	if err := testRepetition(Rand); err != nil {
		t.Error(err)
	}
}

func TestCryptoRepetition(t *testing.T) {
	if err := testRepetition(CryptoRand); err != nil {
		t.Error(err)
	}
}

func testThreadSafety(f func(int) string) {
	const (
		n      = 1000
		length = 15
	)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			f(length)
		}()
	}

	wg.Wait()
}

func TestPseudoThreadSafety(t *testing.T) {
	testThreadSafety(Rand)
}

func TestCryptoThreadSafety(t *testing.T) {
	testThreadSafety(CryptoRand)
}
