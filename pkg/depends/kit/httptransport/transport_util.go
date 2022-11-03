package httptransport

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func tryCatch(fn func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%+v", e)
		}
	}()

	fn()
	return nil
}

func outputln(format string, args ...interface{}) {
	fmt.Printf(`[Kit] `+format+"\n", args...)
}

func methodColor(method string) func(f string, args ...interface{}) string {
	switch method {
	case http.MethodGet:
		return color.BlueString
	case http.MethodPost:
		return color.GreenString
	case http.MethodPut:
		return color.YellowString
	case http.MethodDelete:
		return color.RedString
	default:
		return color.WhiteString
	}
}
