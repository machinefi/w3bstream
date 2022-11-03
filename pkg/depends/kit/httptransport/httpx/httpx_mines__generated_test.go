// This is a generated source file. DO NOT EDIT
// Source: httpx_test/httpx_mines__generated_test.go

package httpx_test

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

func ExampleApplicationOgg() {
	m := httpx.NewApplicationOgg()
	fmt.Println(m.ContentType())
	// Output:
	// application/ogg
}

func ExampleAudioMidi() {
	m := httpx.NewAudioMidi()
	fmt.Println(m.ContentType())
	// Output:
	// audio/midi
}

func ExampleAudioMp3() {
	m := httpx.NewAudioMp3()
	fmt.Println(m.ContentType())
	// Output:
	// audio/mpeg
}

func ExampleAudioOgg() {
	m := httpx.NewAudioOgg()
	fmt.Println(m.ContentType())
	// Output:
	// audio/ogg
}

func ExampleAudioWave() {
	m := httpx.NewAudioWave()
	fmt.Println(m.ContentType())
	// Output:
	// audio/wav
}

func ExampleAudioWebm() {
	m := httpx.NewAudioWebm()
	fmt.Println(m.ContentType())
	// Output:
	// audio/webm
}

func ExampleCSS() {
	m := httpx.NewCSS()
	fmt.Println(m.ContentType())
	// Output:
	// text/css
}

func ExampleHTML() {
	m := httpx.NewHTML()
	fmt.Println(m.ContentType())
	// Output:
	// text/html
}

func ExampleImageBmp() {
	m := httpx.NewImageBmp()
	fmt.Println(m.ContentType())
	// Output:
	// image/bmp
}

func ExampleImageGIF() {
	m := httpx.NewImageGIF()
	fmt.Println(m.ContentType())
	// Output:
	// image/gif
}

func ExampleImageJPEG() {
	m := httpx.NewImageJPEG()
	fmt.Println(m.ContentType())
	// Output:
	// image/jpeg
}

func ExampleImagePNG() {
	m := httpx.NewImagePNG()
	fmt.Println(m.ContentType())
	// Output:
	// image/png
}

func ExampleImageSVG() {
	m := httpx.NewImageSVG()
	fmt.Println(m.ContentType())
	// Output:
	// image/svg+xml
}

func ExampleImageWebp() {
	m := httpx.NewImageWebp()
	fmt.Println(m.ContentType())
	// Output:
	// image/webp
}

func ExamplePlain() {
	m := httpx.NewPlain()
	fmt.Println(m.ContentType())
	// Output:
	// text/plain
}

func ExampleVideoOgg() {
	m := httpx.NewVideoOgg()
	fmt.Println(m.ContentType())
	// Output:
	// video/ogg
}

func ExampleVideoWebm() {
	m := httpx.NewVideoWebm()
	fmt.Println(m.ContentType())
	// Output:
	// video/webm
}
