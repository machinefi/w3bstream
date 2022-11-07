package transformer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"strconv"

	pkgerr "github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func init() { DefaultFactory.Register(&Multipart{}) }

// Multipart for multipart/form-data
type Multipart struct{ *FlattenParams }

func (Multipart) Names() []string {
	return []string{httpx.MIME_MULTIPART_FORMDAT, "multipart", "form-data"}
}

func (Multipart) NamedByTag() string { return "name" }

func (t *Multipart) String() string { return httpx.MIME_MULTIPART_FORMDAT }

func (Multipart) New(ctx context.Context, typ typesx.Type) (Transformer, error) {
	tsf := &Multipart{}

	typ = typesx.DeRef(typ)
	if typ.Kind() != reflect.Struct {
		return nil, pkgerr.Errorf(
			"content transformer `%s` should be used for struct type",
			tsf,
		)
	}

	tsf.FlattenParams = &FlattenParams{}

	if err := tsf.FlattenParams.CollectParams(ctx, typ); err != nil {
		return nil, err
	}
	return tsf, nil
}

func (t *Multipart) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	writer := multipart.NewWriter(w)

	httpx.MaybeWriteHeader(ctx, w, writer.FormDataContentType(), nil)

	errs := vldterr.NewErrorSet()

	for i := range t.Params {
		p := t.Params[i]
		fv := p.FieldValue(rv)

		if p.Tsf != nil {
			st := NewSuper(p.Tsf, &p.Option.CommonOption)

			pw := NewFormPartWriter(
				func(header textproto.MIMEHeader) (io.Writer, error) {
					filename := ""
					if hv := header.Get(httpx.HeaderContentDisposition); hv != "" {
						_, disposition, err := mime.ParseMediaType(hv)
						if err == nil {
							if f, exists := disposition["filename"]; exists {
								filename = fmt.Sprintf("; filename=%s", strconv.Quote(f))
							}
						}
					}
					// always overwrite name
					header.Set(
						httpx.HeaderContentDisposition,
						fmt.Sprintf(
							"form-data; name=%s%s",
							strconv.Quote(p.Name),
							filename,
						),
					)
					return writer.CreatePart(header)
				},
			)

			if err := st.EncodeTo(ctx, pw, fv); err != nil {
				errs.AddErr(err, p.Name)
				continue
			}
		}
	}

	return writer.Close()
}

const dftMaxMemory = 32 << 20 // 32 MB

func (t *Multipart) DecodeFrom(ctx context.Context, r io.Reader, v interface{}, headers ...textproto.MIMEHeader) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	header := MIMEHeader(headers...)
	_, params, err := mime.ParseMediaType(header.Get(httpx.HeaderContentType))
	if err != nil {
		return err
	}

	reader := multipart.NewReader(r, params["boundary"])
	form, err := reader.ReadForm(dftMaxMemory)
	if err != nil {
		return err
	}

	errs := vldterr.NewErrorSet()

	for i := range t.Params {
		p := t.Params[i]
		if p.Tsf == nil {
			continue
		}

		tsf := NewSuper(p.Tsf, &p.Option.CommonOption)

		if files, ok := form.File[p.Name]; ok {
			readers := NewFileHeaderReaders(files)
			if err := tsf.DecodeFrom(ctx, readers, p.FieldValue(rv).Addr()); err != nil {
				errs.AddErr(err, p.Name)
			}
			continue
		}
		if fvs, ok := form.Value[p.Name]; ok {
			readers := NewStringReaders(fvs)
			if err := tsf.DecodeFrom(ctx, readers, p.FieldValue(rv).Addr()); err != nil {
				errs.AddErr(err, p.Name)
			}
		}
	}

	return nil
}

func MustNewFileHeader(fieldName string, filename string, r io.Reader) *multipart.FileHeader {
	header, err := NewFileHeader(fieldName, filename, r)
	if err != nil {
		panic(err)
	}
	return header
}

func NewFileHeader(fieldName string, filename string, r io.Reader) (*multipart.FileHeader, error) {
	buffer := bytes.NewBuffer(nil)
	mpw := multipart.NewWriter(buffer)

	part, err := mpw.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, r); err != nil {
		return nil, err
	}
	mpw.Close()

	reader := multipart.NewReader(buffer, mpw.Boundary())
	form, err := reader.ReadForm(int64(buffer.Len()))
	if err != nil {
		return nil, err
	}

	return form.File[fieldName][0], nil
}

type PartWriterCreator func(header textproto.MIMEHeader) (io.Writer, error)

func NewFormPartWriter(creator PartWriterCreator) *FormPartWriter {
	return &FormPartWriter{
		creator: creator,
		header:  http.Header{},
	}
}

type FormPartWriter struct {
	creator PartWriterCreator
	wr      io.Writer
	header  http.Header
}

func (w *FormPartWriter) NextWriter() io.Writer { return NewFormPartWriter(w.creator) }

func (w *FormPartWriter) Header() http.Header { return w.header }

func (w *FormPartWriter) Write(p []byte) (n int, err error) {
	if w.wr == nil {
		w.wr, err = w.creator(textproto.MIMEHeader(w.header))
		if err != nil {
			return -1, err
		}
	}
	return w.wr.Write(p)
}

func NewFileHeaderReaders(headers []*multipart.FileHeader) *StringReaders {
	bs := make([]io.Reader, len(headers))
	for i := range headers {
		bs[i] = &FileHeaderReader{v: headers[i]}
	}

	return &StringReaders{
		readers: bs,
	}
}

type FileHeaderReader struct {
	v      *multipart.FileHeader
	opened multipart.File
}

func (f *FileHeaderReader) Interface() interface{} { return f.v }

func (f *FileHeaderReader) Read(p []byte) (int, error) {
	if f.opened == nil {
		file, err := f.v.Open()
		if err != nil {
			return -1, err
		}
		f.opened = file
	}
	n, err := f.opened.Read(p)
	if err == io.EOF {
		return n, f.opened.Close()
	}
	return n, err
}
