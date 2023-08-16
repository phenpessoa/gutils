// Package unsafex, like unsafe, contains operations that step around the type
// safety of Go programs.
//
// Packages that import unsafex may be non-portable and are not protected by the
// Go 1 compatibility guidelines.
package unsafex

import "unsafe"

// String is a wrapper over unsafe.String, converting b to a string.
//
// Since Go strings are immutable, the bytes passed to String must not be
// modified afterwards.
func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// ByteSlice is a wrapper over unsafe.Slice, converting s to a []byte.
//
// Since Go strings are immutable, the bytes received from ByteSlice must not be
// modified.
func ByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
