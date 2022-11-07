// Package ptrx used for inline point value
package ptrx

import "time"

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Bool(v bool) *bool { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Int(v int) *int { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Int8(v int8) *int8 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Int16(v int16) *int16 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Int32(v int32) *int32 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Int64(v int64) *int64 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Uint(v uint) *uint { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Uint8(v uint8) *uint8 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Uint16(v uint16) *uint16 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Uint32(v uint32) *uint32 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Uint64(v uint64) *uint64 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Float32(v float32) *float32 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Float64(v float64) *float64 { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Byte(v byte) *byte { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Rune(v rune) *rune { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func String(v string) *string { return &v }

// Deprecated: use Ptr instead. if your go version greater than 1.18
func Duration(d time.Duration) *time.Duration { return &d }

func Ptr[V any](v V) *V { return &v }
