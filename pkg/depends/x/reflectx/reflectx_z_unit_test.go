package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type Foo struct {
	Field int `name:"fieldNameTag,default='0'" json:"fieldJsonTag,omitempty"`
}

func TestNatureType(t *testing.T) {
	var v, vp = Foo{}, &Foo{}

	NewWithT(t).Expect(NatureType(v).String()).To(Equal("reflectx_test.Foo"))
	NewWithT(t).Expect(NatureType(vp).String()).To(Equal("reflectx_test.Foo"))
	NewWithT(t).Expect(NatureType(&vp).String()).To(Equal("reflectx_test.Foo"))
	NewWithT(t).Expect(NatureType(reflect.TypeOf(v)).String()).To(Equal("reflectx_test.Foo"))
	NewWithT(t).Expect(NatureType(reflect.TypeOf(vp)).String()).To(Equal("reflectx_test.Foo"))
	NewWithT(t).Expect(NatureType(reflect.TypeOf(&vp)).String()).To(Equal("reflectx_test.Foo"))

}

func TestTagValueAndFlags(t *testing.T) {
	ft, _ := reflect.ValueOf(Foo{}).Type().FieldByName("Field")
	nameTag, _ := ft.Tag.Lookup("name")
	jsonTag, _ := ft.Tag.Lookup("json")

	key, flags := TagValueAndFlags(nameTag)
	NewWithT(t).Expect(key).To(Equal("fieldNameTag"))
	NewWithT(t).Expect(flags).To(Equal(map[string]bool{"default='0'": true}))

	key, flags = TagValueAndFlags(jsonTag)
	NewWithT(t).Expect(key).To(Equal("fieldJsonTag"))
	NewWithT(t).Expect(flags).To(Equal(map[string]bool{"omitempty": true}))
}
