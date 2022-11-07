package statusx

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Error interface {
	StatusErr() *StatusErr
	Error() string
}

type ServiceCode interface {
	ServiceCode() int
}

func IsStatusErr(err error) (*StatusErr, bool) {
	if err == nil {
		return nil, false
	}

	if e, ok := err.(Error); ok {
		return e.StatusErr(), ok
	}

	se, ok := err.(*StatusErr)
	return se, ok
}

func FromErr(err error) *StatusErr {
	if err == nil {
		return nil
	}
	if se, ok := IsStatusErr(err); ok {
		return se
	}
	return NewUnknownErr().WithDesc(err.Error())
}

func Wrap(err error, code int, key string, msgs ...string) *StatusErr {
	if err == nil {
		return nil
	}

	if len(strconv.Itoa(code)) == 3 {
		code = code * 1e6
	}

	msg := key

	if len(msgs) > 0 {
		msg = msgs[0]
	}

	desc := ""

	if len(msgs) > 1 {
		desc = strings.Join(msgs[1:], "\n")
	} else {
		desc = err.Error()
	}

	// err = errors.WithMessage(err, "asdfasdfasdfasdfass")
	s := &StatusErr{
		Key:   key,
		Code:  code,
		Msg:   msg,
		Desc:  desc,
		error: errors.WithStack(err),
	}

	return s
}

func NewUnknownErr() *StatusErr {
	return NewStatusErr("UnknownError", http.StatusInternalServerError*1e6, "unknown error")
}

func NewStatusErr(key string, code int, msg string) *StatusErr {
	return &StatusErr{
		Key:  key,
		Code: code,
		Msg:  msg,
	}
}

type StatusErr struct {
	Key       string      `json:"key"       xml:"key"`       // key of err
	Code      int         `json:"code"      xml:"code"`      // unique err code
	Msg       string      `json:"msg"       xml:"msg"`       // msg of err
	Desc      string      `json:"desc"      xml:"desc"`      // desc of err
	CanBeTalk bool        `json:"canBeTalk" xml:"canBeTalk"` // can be task error; for client to should error msg to end user
	ID        string      `json:"id"        xml:"id"`        // request ID or other request context
	Sources   []string    `json:"sources"   xml:"sources"`   // error tracing
	Fields    ErrorFields `json:"fields"    xml:"fields"`    // error in where fields
	error     error
}

// @err[UnknownError][500000000][unknown error]
var regexpStatusErrSummary = regexp.MustCompile(`@StatusErr\[(.+)\]\[(.+)\]\[(.+)\](!)?`)

func ParseStatusErrSummary(s string) (*StatusErr, error) {
	if !regexpStatusErrSummary.Match([]byte(s)) {
		return nil, fmt.Errorf("unsupported status err summary: %s", s)
	}

	matched := regexpStatusErrSummary.FindStringSubmatch(s)

	code, _ := strconv.ParseInt(matched[2], 10, 64)

	return &StatusErr{
		Key:       matched[1],
		Code:      int(code),
		Msg:       matched[3],
		CanBeTalk: matched[4] != "",
	}, nil
}

func (se *StatusErr) Summary() string {
	s := fmt.Sprintf(
		`@StatusErr[%s][%d][%s]`,
		se.Key,
		se.Code,
		se.Msg,
	)

	if se.CanBeTalk {
		return s + "!"
	}
	return s
}

func (se *StatusErr) Is(err error) bool {
	e := FromErr(err)
	if se == nil || e == nil {
		return false
	}
	return e.Key == se.Key && e.Code == se.Code
}

func StatusCodeFromCode(code int) int {
	strCode := fmt.Sprintf("%d", code)
	if len(strCode) < 3 {
		return 0
	}
	statusCode, _ := strconv.Atoi(strCode[:3])
	return statusCode
}

func (se *StatusErr) StatusCode() int {
	return StatusCodeFromCode(se.Code)
}

func (se *StatusErr) Error() string {
	s := fmt.Sprintf(
		"[%s]%s%s",
		strings.Join(se.Sources, ","),
		se.Summary(),
		se.Fields,
	)

	if se.Desc != "" {
		s += " " + se.Desc
	}

	return s
}

func (se StatusErr) WithMsg(msg string) *StatusErr {
	se.Msg = msg
	return &se
}

func (se StatusErr) WithDesc(desc string) *StatusErr {
	se.Desc = desc
	return &se
}

func (se StatusErr) WithID(id string) *StatusErr {
	se.ID = id
	return &se
}

func (se StatusErr) AppendSource(sourceName string) *StatusErr {
	length := len(se.Sources)
	if length == 0 || se.Sources[length-1] != sourceName {
		se.Sources = append(se.Sources, sourceName)
	}
	return &se
}

func (se StatusErr) EnableErrTalk() *StatusErr {
	se.CanBeTalk = true
	return &se
}

func (se StatusErr) DisableErrTalk() *StatusErr {
	se.CanBeTalk = false
	return &se
}

func (se StatusErr) AppendErrorField(in string, field string, msg string) *StatusErr {
	se.Fields = append(se.Fields, NewErrorField(in, field, msg))
	return &se
}

func (se StatusErr) AppendErrorFields(errorFields ...*ErrorField) *StatusErr {
	se.Fields = append(se.Fields, errorFields...)
	return &se
}

func NewErrorField(in string, field string, msg string) *ErrorField {
	return &ErrorField{
		In:    in,
		Field: field,
		Msg:   msg,
	}
}

type ErrorField struct {
	Field string `json:"field" xml:"field"` // Field path: prop.slice[2].a
	Msg   string `json:"msg"   xml:"msg"`   // Msg message
	In    string `json:"in"    xml:"in"`    // In location eq. body, query, header, path, formData
}

func (s ErrorField) String() string {
	return s.Field + " in " + s.In + " - " + s.Msg
}

type ErrorFields []*ErrorField

func (fs ErrorFields) String() string {
	if len(fs) == 0 {
		return ""
	}

	sort.Sort(fs)

	buf := &bytes.Buffer{}
	buf.WriteString("<")
	for i, f := range fs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(f.String())
	}
	buf.WriteString(">")
	return buf.String()
}

func (fs ErrorFields) Len() int {
	return len(fs)
}

func (fs ErrorFields) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs ErrorFields) Less(i, j int) bool {
	return fs[i].Field < fs[j].Field
}
