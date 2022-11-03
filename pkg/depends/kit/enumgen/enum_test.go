package enumgen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

var (
	g      *enumgen.Generator
	sample *enumgen.Enum
	scheme *enumgen.Enum
	policy *enumgen.Enum
	f      *codegen.File
	PkgID  string
)

func init() {
	cwd, _ := os.Getwd()
	pkg, err := pkgx.LoadFrom(filepath.Join(cwd, "./__examples__"))
	if err != nil {
		panic(err)
	}
	PkgID = pkg.ID
	g = enumgen.New(pkg)
	g.Scan("Sample", "Scheme", "PullPolicy")

	if sample = enumgen.GetEnumByName(g, "Sample"); sample == nil {
		panic(nil)
	}
	if scheme = enumgen.GetEnumByName(g, "Scheme"); scheme == nil {
		panic(nil)
	}
	// should be nil, just generate IntStringerEnum
	if policy = enumgen.GetEnumByName(g, "PullPolicy"); policy != nil {
		panic(nil)
	}
	g.Output(cwd)
	// fake f for testing enum member functions
	f = codegen.NewFile(pkg.Name, "mock.go")
}

func ExampleEnum_ConstName() {
	fmt.Println(string(sample.ConstName("XXX").Bytes()))
	fmt.Println(string(sample.ConstName("ABC").Bytes()))
	fmt.Println(string(scheme.ConstName("HTTP").Bytes()))
	fmt.Println(string(scheme.ConstName("HTTPS").Bytes()))
	// Output:
	// SAMPLE__XXX
	// SAMPLE__ABC
	// SCHEME__HTTP
	// SCHEME__HTTPS
}

func ExampleEnum_Errors() {
	fmt.Println(string(sample.Errors(f).Bytes()))
	fmt.Println(string(scheme.Errors(f).Bytes()))
	// Output:
	// var InvalidSample = errors.New("invalid Sample type")
	// var InvalidScheme = errors.New("invalid Scheme type")
}

func ExampleEnum_StringParser() {
	fmt.Println(string(sample.StringParser(f).Bytes()))
	// Output:
	// func ParseSampleFromString(s string) (Sample, error) {
	// switch s {
	// default:
	// return SAMPLE_UNKNOWN, InvalidSample
	// case "":
	// return SAMPLE_UNKNOWN, nil
	// case "XXX":
	// return SAMPLE__XXX, nil
	// case "YYY":
	// return SAMPLE__YYY, nil
	// case "ZZZ":
	// return SAMPLE__ZZZ, nil
	// }
	// }
}

func ExampleEnum_LabelParser() {
	fmt.Println(string(sample.LabelParser(f).Bytes()))
	// Output:
	// func ParseSampleFromLabel(s string) (Sample, error) {
	// switch s {
	// default:
	// return SAMPLE_UNKNOWN, InvalidSample
	// case "":
	// return SAMPLE_UNKNOWN, nil
	// case "样例XXX":
	// return SAMPLE__XXX, nil
	// case "样例YYY":
	// return SAMPLE__YYY, nil
	// case "样例ZZZ":
	// return SAMPLE__ZZZ, nil
	// }
	// }
}

func ExampleEnum_Stringer() {
	fmt.Println(string(sample.Stringer(f).Bytes()))
	// Output:
	// func (v Sample) String() string {
	// switch v {
	// default:
	// return "UNKNOWN"
	// case SAMPLE_UNKNOWN:
	// return ""
	// case SAMPLE__XXX:
	// return "XXX"
	// case SAMPLE__YYY:
	// return "YYY"
	// case SAMPLE__ZZZ:
	// return "ZZZ"
	// }
	// }
}

func ExampleEnum_Integer() {
	fmt.Println(string(sample.Integer(f).Bytes()))
	// Output:
	// func (v Sample) Int() int {
	// return int(v)
	// }
}

func ExampleEnum_Labeler() {
	fmt.Println(string(sample.Labeler(f).Bytes()))
	// Output:
	// func (v Sample) Label() string {
	// switch v {
	// default:
	// return "UNKNOWN"
	// case SAMPLE_UNKNOWN:
	// return ""
	// case SAMPLE__XXX:
	// return "样例XXX"
	// case SAMPLE__YYY:
	// return "样例YYY"
	// case SAMPLE__ZZZ:
	// return "样例ZZZ"
	// }
	// }
}

func TestEnum_TypeName(t *testing.T) {
	NewWithT(t).Expect(sample.TypeName(f).Bytes()).To(Equal([]byte(
		`func (v Sample) TypeName() string {
return "` + PkgID + `.Sample"
}`)))
}

func ExampleEnum_ConstValues() {
	fmt.Println(string(sample.ConstValues(f).Bytes()))
	// Output:
	// func (v Sample) ConstValues() []enum.IntStringerEnum {
	// return []enum.IntStringerEnum{SAMPLE__XXX, SAMPLE__YYY, SAMPLE__ZZZ}
	// }
}

func ExampleEnum_TextMarshaler() {
	fmt.Println(string(sample.TextMarshaler(f).Bytes()))
	// Output:
	// func (v Sample) MarshalText() ([]byte, error) {
	// s := v.String()
	// if s == "UNKNOWN" {
	// return nil, InvalidSample
	// }
	// return []byte(s), nil
	// }
}

func ExampleEnum_TextUnmarshaler() {
	fmt.Println(string(sample.TextUnmarshaler(f).Bytes()))
	// Output:
	// func (v *Sample) UnmarshalText(data []byte) error {
	// s := string(bytes.ToUpper(data))
	// val, err := ParseSampleFromString(s)
	// if err != nil {
	// return err
	// }
	// *(v) = val
	// return nil
	// }
}

func ExampleEnum_Scanner() {
	fmt.Println(string(sample.Scanner(f).Bytes()))
	// Output:
	// func (v *Sample) Scan(src interface{}) error {
	// offset := 0
	// o, ok := interface{}(v).(enum.ValueOffset)
	// if ok {
	// offset = o.Offset()
	// }
	// i, err := enum.ScanIntEnumStringer(src, offset)
	// if err != nil {
	// return err
	// }
	// *(v) = Sample(i)
	// return nil
	// }
}

func ExampleEnum_Valuer() {
	fmt.Println(string(sample.Valuer(f).Bytes()))
	// Output:
	// func (v Sample) Value() (driver.Value, error) {
	// offset := 0
	// o, ok := interface{}(v).(enum.ValueOffset)
	// if ok {
	// offset = o.Offset()
	// }
	// return int64(v) + int64(offset), nil
	// }
}
