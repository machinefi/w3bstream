// This is a generated source file. DO NOT EDIT
// Source: must/must_generated.go

package must

import "log"

func Byte(v byte, err error) byte {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func ByteOK(v byte, ok bool) byte {
	if !ok {
		log.Panic("Byte not ok")
	}
	return v
}

func Bytes(v []byte, err error) []byte {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func BytesOK(v []byte, ok bool) []byte {
	if !ok {
		log.Panic("Bytes not ok")
	}
	return v
}

func String(v string, err error) string {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func StringOK(v string, ok bool) string {
	if !ok {
		log.Panic("String not ok")
	}
	return v
}

func Strings(v []string, err error) []string {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func StringsOK(v []string, ok bool) []string {
	if !ok {
		log.Panic("Strings not ok")
	}
	return v
}

func Int(v int, err error) int {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func IntOK(v int, ok bool) int {
	if !ok {
		log.Panic("Int not ok")
	}
	return v
}

func Int8(v int8, err error) int8 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Int8OK(v int8, ok bool) int8 {
	if !ok {
		log.Panic("Int8 not ok")
	}
	return v
}

func Int16(v int16, err error) int16 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Int16OK(v int16, ok bool) int16 {
	if !ok {
		log.Panic("Int16 not ok")
	}
	return v
}

func Int32(v int32, err error) int32 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Int32OK(v int32, ok bool) int32 {
	if !ok {
		log.Panic("Int32 not ok")
	}
	return v
}

func Int64(v int64, err error) int64 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Int64OK(v int64, ok bool) int64 {
	if !ok {
		log.Panic("Int64 not ok")
	}
	return v
}

func Uint8(v uint8, err error) uint8 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Uint8OK(v uint8, ok bool) uint8 {
	if !ok {
		log.Panic("Uint8 not ok")
	}
	return v
}

func Uint16(v uint16, err error) uint16 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Uint16OK(v uint16, ok bool) uint16 {
	if !ok {
		log.Panic("Uint16 not ok")
	}
	return v
}

func Uint32(v uint32, err error) uint32 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Uint32OK(v uint32, ok bool) uint32 {
	if !ok {
		log.Panic("Uint32 not ok")
	}
	return v
}

func Uint64(v uint64, err error) uint64 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Uint64OK(v uint64, ok bool) uint64 {
	if !ok {
		log.Panic("Uint64 not ok")
	}
	return v
}

func Rune(v rune, err error) rune {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func RuneOK(v rune, ok bool) rune {
	if !ok {
		log.Panic("Rune not ok")
	}
	return v
}

func Float32(v float32, err error) float32 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Float32OK(v float32, ok bool) float32 {
	if !ok {
		log.Panic("Float32 not ok")
	}
	return v
}

func Float64(v float64, err error) float64 {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Float64OK(v float64, ok bool) float64 {
	if !ok {
		log.Panic("Float64 not ok")
	}
	return v
}
