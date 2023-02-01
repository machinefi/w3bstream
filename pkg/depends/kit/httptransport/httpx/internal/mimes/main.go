package main

import (
	"path/filepath"
	"runtime"
	"sort"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
var mimes = []struct{ typ, name string }{
	{"text/plain", "Plain"},
	{"text/css", "CSS"},
	{"text/html", "HTML"},
	{"image/gif", "ImageGIF"},
	{"image/jpeg", "ImageJPEG"},
	{"image/png", "ImagePNG"},
	{"image/bmp", "ImageBmp"},
	{"image/webp", "ImageWebp"},
	{"image/svg+xml", "ImageSVG"},
	{"audio/wav", "AudioWave"},
	{"audio/midi", "AudioMidi"},
	{"audio/webm", "AudioWebm"},
	{"video/webm", "VideoWebm"},
	{"audio/ogg", "AudioOgg"},
	{"video/ogg", "VideoOgg"},
	{"application/ogg", "ApplicationOgg"},
	{"application/x-protobuf", "ApplicationProtobuf"},
	{"audio/mpeg", "AudioMp3"},
}

func main() {
	sort.Slice(mimes, func(i, j int) bool {
		return mimes[i].name < mimes[j].name
	})
	{
		f := g.NewFile("httpx", g.GenerateFileSuffix("./httpx_mines.go"))
		keys := make([]string, 0)

		sort.Strings(keys)

		for _, m := range mimes {
			t := g.Type(m.name)
			f.WriteSnippet(
				f.Expr(`func New?() *? {return &?{}}`, t, t, t),
				f.Expr(`type ? struct {`+f.Use("bytes", "Buffer")+`}`, t),
				f.Expr(`func (?) ContentType() string {return ?}`, t, f.Value(m.typ)),
			)
		}
		if _, err := f.Write(); err != nil {
			panic(err)
		}

	}

	{
		f := g.NewFile("httpx_test", g.GenerateFileSuffix("./httpx_mines_test.go"))

		for _, m := range mimes {
			t := g.Type(m.name)
			fn := "New" + m.name
			f.WriteSnippet(
				f.Expr(`func Example?() {
m := `+f.Use(pkg, fn)+`()
`+f.Use("fmt", "Println")+`(m.ContentType())
// Output:
// `+m.typ+`
}`, t, t),
			)
		}

		if _, err := f.Write(); err != nil {
			panic(err)
		}
	}
}

var pkg string

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(current), "../../../httpx")
	pkg = must.String(pkgx.PkgIdByPath(dir))
}
