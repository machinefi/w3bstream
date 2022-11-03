package statusx_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusxgen"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/statusxgen/__examples__"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func init() {
	cwd, _ := os.Getwd()
	pkg, _ := pkgx.LoadFrom(filepath.Join(cwd, "../statusxgen/__examples__"))

	g := statusxgen.New(pkg)

	g.Scan("StatusError")
	g.Output(cwd)
}

func TestStatusErr(t *testing.T) {
	var (
		unknownSeStr  = "@StatusErr[UnknownError][500000000][unknown error]"
		unauthedSeStr = "@StatusErr[Unauthorized][401999001][Unauthorized]!"
		internalSeStr = "@StatusErr[InternalServerError][500999001][InternalServerError 内部错误]"
		summary       = statusx.NewUnknownErr().Summary()
		se, err       = statusx.ParseStatusErrSummary(summary)
	)

	g := NewWithT(t)

	g.Expect(summary).To(Equal(unknownSeStr))
	g.Expect(err).To(BeNil())
	g.Expect(se).To(Equal(statusx.NewUnknownErr()))
	g.Expect(Unauthorized.StatusErr().Summary()).To(Equal(unauthedSeStr))
	g.Expect(InternalServerError.StatusErr().Summary()).To(Equal(internalSeStr))
	g.Expect(Unauthorized.StatusCode()).To(Equal(401))
	g.Expect(Unauthorized.StatusErr().StatusCode()).To(Equal(401))

	g.Expect(errors.Is(Unauthorized, Unauthorized)).To(BeTrue())
	g.Expect(errors.Is(Unauthorized.StatusErr(), Unauthorized)).To(BeTrue())
	g.Expect(errors.Is(Unauthorized.StatusErr(), Unauthorized.StatusErr())).To(BeTrue())
}

func ExampleStatusErr() {
	fmt.Println(Unauthorized)
	fmt.Println(statusx.FromErr(nil))
	fmt.Println(statusx.FromErr(fmt.Errorf("unknown")))
	fmt.Println(Unauthorized.StatusErr().WithMsg("msg overwrite"))
	fmt.Println(Unauthorized.StatusErr().WithDesc("desc overwrite"))
	fmt.Println(Unauthorized.StatusErr().DisableErrTalk().EnableErrTalk())
	fmt.Println(Unauthorized.StatusErr().WithID("111"))
	fmt.Println(Unauthorized.StatusErr().AppendSource("service-abc"))
	fmt.Println(Unauthorized.StatusErr().AppendErrorField("header", "Authorization", "missing"))
	fmt.Println(Unauthorized.StatusErr().AppendErrorFields(
		statusx.NewErrorField("query", "key", "missing"),
		statusx.NewErrorField("header", "Authorization", "missing"),
	))
	// Output:
	// []@StatusErr[Unauthorized][401999001][Unauthorized]!
	// <nil>
	// []@StatusErr[UnknownError][500000000][unknown error] unknown
	// []@StatusErr[Unauthorized][401999001][msg overwrite]!
	// []@StatusErr[Unauthorized][401999001][Unauthorized]! desc overwrite
	// []@StatusErr[Unauthorized][401999001][Unauthorized]!
	// []@StatusErr[Unauthorized][401999001][Unauthorized]!
	// [service-abc]@StatusErr[Unauthorized][401999001][Unauthorized]!
	// []@StatusErr[Unauthorized][401999001][Unauthorized]!<Authorization in header - missing>
	// []@StatusErr[Unauthorized][401999001][Unauthorized]!<Authorization in header - missing, key in query - missing>
}
