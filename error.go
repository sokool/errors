package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"runtime"
	"strconv"
	"strings"
)

type Error struct {
	wrapped error
	message string
	tag     string
	code    string
	call    int
	id      uint32
	trace   struct {
		file string
		line int
		funk string
	}
}

func New(message string, args ...any) *Error {
	return Trace(1, message, args...)
}

func Trace(deep int, message string, args ...any) *Error {
	e := Error{message: strings.TrimSpace(fmt.Errorf(message, args...).Error())}
	var m = fmt.Errorf(message, args...).Error()
	if i := strings.Index(m, ":"); i > 0 && !strings.Contains(m[:i], " ") {
		e.tag, e.message = m[:i], m[i+1:]
		if k := strings.Index(e.tag, "#"); k >= 0 {
			e.code, e.tag = e.tag[k+1:], e.tag[:k]
		}
		//fmt.Printf("tag:%s, code:%s, message:%s, m:%s\n", e.tag, e.code, e.message, m)
	}
	if p, f, l, ok := runtime.Caller(deep + 1); ok {
		e.trace.funk, e.trace.file, e.trace.line = runtime.FuncForPC(p).Name(), f, l
	}

	//stack := make([]uintptr, 50)
	//runtime.Callers(deep+2, stack[:])
	//x := runtime.CallersFrames(stack)
	//for {
	//	f, ok := x.Next()
	//	if !ok {
	//		break
	//	}
	//	fmt.Printf("%s:%d -> %s\n", f.File, f.Line, f.Function)
	//}

	var s, k = len(message), 0
	for i := range message {
		if i+2 <= s && message[i:i+2] == "%w" {
			e.wrapped = args[k].(error)
			break
		}
		if message[i] == '%' {
			k++
		}
	}

	if e.id == 0 {
		t := fmt.Sprintf("%s%s%s", e.code, e.tag, fmt.Sprintf(message, args...))
		h := fnv.New32a()
		h.Write([]byte(t))
		e.id = h.Sum32()
	}

	return &e
}

func (e *Error) ID() uint32 {
	return e.id
}

func (e *Error) Exp() {
}

func (e *Error) Tag() string {
	return e.tag
}

func (e *Error) Func() string {
	return e.trace.funk
}

func (e *Error) Line() int {
	return e.trace.line
}

func (e *Error) File() string {
	return e.trace.file
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) CodeNumber() int {
	if n, _ := strconv.Atoi(e.code); n > 0 {
		return n
	}
	return -1
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {

	c := e.Code()
	t := e.Tag()
	m := e.Message()
	if c != "" && t != "" {
		m = fmt.Sprintf("%s#%s:%s", t, c, m)
	} else if c != "" {
		m = fmt.Sprintf("#%s:%s", c, m)
	} else if t != "" {
		m = fmt.Sprintf("%s:%s", t, m)
	}
	strings.TrimSpace(m)
	return strings.TrimSpace(m)
	//email#h1: invalid hostname

	//return fmt.Sprintf("%s %s:%d", e.Message(), e.trace.file, e.trace.line)
	s := fmt.Sprintf("%s#%s:%s", e.tag, e.code, e.Message())
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

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Unwrap() error {
	return e.wrapped
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":      e.id,
		"message": e.message,
		"tag":     e.tag,
		"code":    e.code,
		"file":    e.trace.file,
		"line":    e.trace.line,
		"func":    e.trace.funk,
	})
}

func Extract(from error) (o []error) {
	o = append(o, from)
	for {
		var err error
		if err = errors.Unwrap(from); err == nil {
			return
		}
		o = append(o, err)
		from = err
	}
}

func First(from error) *Error {
	var e *Error
	if errors.As(from, &e) {
		return e
	}
	return nil
}
