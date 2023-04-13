package definition_scanner

import (
	"mime/multipart"
	"time"

	time2 "github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/testdata/definition_scanner/time"
)

type Interface interface{}

type Binary *multipart.FileHeader
type String string
type Bool bool
type Float float32
type Double float64
type Int int
type Uint uint

// 日期
type Time time.Time

func (Time) OpenAPISchemaType() []string { return []string{"string"} }
func (Time) OpenAPISchemaFormat() string { return "date-time" }

// openapi:type boolean
type FakeBool int

type Map map[string]String

type ArrayString [2]string
type SliceString []string
type SliceNamed []String

type TimeAlias = time.Time

type Struct struct {
	// name
	Name      *string    `json:"name" validate:"@string[2,]"`
	ID        **string   `json:"id,omitempty" default:"1" validate:"@string/\\d+/"`
	Enum      Enum       `json:"enum" validate:"@string{ONE}"`
	CreatedAt time2.Time `json:"createdAt,omitempty"`
	Slice     []float64  `json:"slice" validate:"@slice<@float64<7,5>>[1,3]"`
	Map       map[string]map[string]struct {
		ID int `json:"id" validate:"@int[0,10]"`
	} `json:"map,omitempty" validate:"@map<,@map<,@struct>>[0,3]"`
}

type Part struct {
	Name string `json:",omitempty" validate:"@string[2,]"`
	skip string
	Skip string `json:"-"`
}

type PartConflict struct {
	Name string `json:"name" validate:"@string[0,)"`
}

// Composed
type Composed struct {
	Part
}

type NamedComposed struct {
	Part `json:"part"`
}

type InvalidComposed struct {
	Part
	PartConflict
}

type Node struct {
	Type     string  `json:"type"`
	Children []*Node `json:"children"`
}
