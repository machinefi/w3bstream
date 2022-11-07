package scanner

import (
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type ScanIterator interface {
	// New a ptr value for scan
	New() interface{}
	// Next for receive scanned value
	Next(v interface{}) error
}

func ScanIteratorFor(v interface{}) (ScanIterator, error) {
	switch x := v.(type) {
	case ScanIterator:
		return x, nil
	default:
		t := reflectx.DeRef(reflect.TypeOf(v))
		if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8 {
			return &SliceScanIterator{
				t.Elem(), reflectx.Indirect(reflect.ValueOf(v)),
			}, nil
		}
		return &SingleScanIterator{v: v}, nil
	}
}

type SliceScanIterator struct {
	t reflect.Type
	v reflect.Value
}

func (s *SliceScanIterator) New() interface{} {
	return reflectx.New(s.t).Addr().Interface()
}

func (s *SliceScanIterator) Next(v interface{}) error {
	s.v.Set(reflect.Append(s.v, reflect.ValueOf(v).Elem()))
	return nil
}

type SingleScanIterator struct {
	v          interface{}
	hasResults bool
}

func (s *SingleScanIterator) New() interface{} {
	return s.v
}

func (s *SingleScanIterator) Next(v interface{}) error {
	s.hasResults = true
	return nil
}

func (s *SingleScanIterator) HasRecord() bool { return s.hasResults }
