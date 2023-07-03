package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	message, name, code string
	prev                error
}

func New(message, name, code string, caller ...int) *Error {
	if len(caller) == 0 {
		caller = append(caller, 1)
	}
	var err Error
	message, name, code = strings.TrimSpace(message), strings.TrimSpace(name), strings.TrimSpace(code)
	if name == "" && code == "" {
		if p, f, l, ok := runtime.Caller(caller[0]); ok {
			name, code = runtime.FuncForPC(p).Name(), fmt.Sprintf("%s@L%d", f[strings.LastIndex(f, "/")+1:], l)
		}
	}
	err.message, err.name, err.code = message, name, code
	return &err
}

func Errorf(format string, args ...any) *Error {
	var m, n, c = fmt.Errorf(format, args...).Error(), "", ""
	if i := strings.Index(m, ":"); i > 0 && !strings.Contains(m[:i], " ") {
		n, m = m[:i], m[i+1:]
		if k := strings.Index(n, "#"); k >= 0 {
			c, n = n[k+1:], n[:k]
		}
	}

	e := New(m, n, c, 2)

	// find wrapped error and store it in domain.Error
	var s = len(format)
	var k int
	for i := range format {
		if i+2 <= s && format[i:i+2] == "%w" {
			e.prev = args[k].(error)
			break
		}
		if format[i] == '%' {
			k++
		}
	}

	return e
}

func Read(from error) *Error {
	var e *Error
	if errors.As(from, &e) {
		return e
	}
	return nil
}

func (e *Error) Name() string {
	return e.name
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	s := fmt.Sprintf("%s#%s: %s", e.name, e.code, e.message)
	if s[:3] == "#: " {
		s = s[3:]
	}
	if n := len(s); n > 2 && s[n-2:] == ": " {
		s = s[:n-2]
	}
	if n := len(s); n > 1 && s[n-1:] == "#" {
		s = s[:n-1] + ":"
	}
	return strings.Replace(s, "#:", ":", -1)
}

func (e *Error) Unwrap() error {
	return e.prev
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{"message": e.message, "name": e.name, "code": e.code})
}

func Trace(err error) (o []error) {
	o = append(o, err)
	for {
		//s1 := err.Error()
		var err2 error
		if err2 = errors.Unwrap(err); err2 == nil {
			//fmt.Printf("%T: %s\n", err2, s1)
			return
		}

		//s2 := err2.Error()
		//s := strings.Replace(s1, s2, "", -1)
		//fmt.Printf("%T: %s\n", err, s)
		o = append(o, err2)
		err = err2
	}
}
