package errors_test

import (
	"fmt"
	"testing"

	"github.com/sokool/errors"
)

func TestNew(t *testing.T) {
	type scenario struct {
		description                            string
		input                                  *errors.Error // given
		tag, code, file, funk, message, output string        // expects
		line                                   int
	}

	cases := []scenario{
		{
			description: "empty",
			input:       errors.New(""),
			file:        "/Users/sokool/Code/Sokool/errors/error_test.go",
			funk:        "github.com/sokool/errors_test.TestNew",
			line:        21,
		},
		{
			description: "just message",
			input:       errors.New("hi there"),
			file:        "/Users/sokool/Code/Sokool/errors/error_test.go",
			funk:        "github.com/sokool/errors_test.TestNew",
			line:        28,
			output:      "hi there",
		},
		{
			description: "message with arguments",
			input:       errors.New("hi there %s", "man"),
			output:      "hi there man",
		},
		{
			description: "message with wrapped error",
			input:       errors.New("one then %w", fmt.Errorf("two")),
			output:      "one then two",
		},
		{
			description: "only tag name",
			input:       errors.New("test:"),
			tag:         "test",
			output:      "test:",
		},
		{
			description: "only code",
			input:       errors.New("#h6b7:"),
			output:      "#h6b7:",
			code:        "h6b7",
		},
		{
			description: "tag,message",
			input:       errors.New("test: hi there"),
			tag:         "test",
			output:      "test: hi there",
		},
		{
			description: "tag,code,message",
			input:       errors.New("email#h1: invalid hostname"),
			tag:         "email",
			code:        "h1",
			//message:     "invalid hostname",
			output: "email#h1: invalid hostname",
		},
		{
			description: "with name,code,message from arguments",
			input:       errors.New("%s#%s: invalid %s", "email", "h45", "username"),
			tag:         "email",
			code:        "h45",
			//message:     "invalid username",
			output: "email#h45: invalid username",
		},
		{
			description: "with code,message",
			input:       errors.New("#e87:failed"),
			code:        "e87",
			//message:     "failed",
			output: "#e87:failed",
		},
		{
			description: "with code,name signs in wrong place",
			input:       errors.New("failed due abc: #triggered"),
			line:        87,
			//file:        "error_test.go:85",
			//message:     "failed due abc: #triggered",
			output: "failed due abc: #triggered",
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Helper()
			if s := c.input.Tag(); s != c.tag {
				t.Fatalf("expected tag `%s`, got `%s`", c.tag, s)
			}
			if s := c.input.Code(); s != c.code {
				t.Fatalf("expected code `%s`, got `%s`", c.code, s)
			}
			if s := c.input.Message(); c.message != "" && s != c.message {
				t.Fatalf("expected message `%s`, got `%s`", c.message, s)
			}
			if s := c.input.Error(); s != c.output {
				t.Fatalf("expected output `%s`, got `%s`", c.output, s)
			}
			if s := c.input.Func(); c.funk != "" && s != c.funk {
				t.Fatalf("expected `%s` func, got `%s`", c.funk, s)
			}
			if s := c.input.File(); c.file != "" && s != c.file {
				t.Fatalf("expected `%s` file, got `%s`", c.file, s)
			}
			if s := c.input.Line(); c.line != 0 && s != c.line {
				t.Fatalf("expected %d file line, got %d", c.line, s)
			}
		})
	}
}

func TestError_Line(t *testing.T) {
	errorf := func(message string, args ...any) error {
		return func(message string, args ...any) error {
			return Foo()
		}(message, args...)
	}
	x := errorf("").(*errors.Error)
	if s := x.Line(); s != 182 {
		t.Fatalf("got %d expected line 136", s)
	}
}

//func TestFirst(t *testing.T) {
//	Foo()
//	fmt.Println()
//	a := errors.New("eloszki")
//	b := fmt.Errorf("second %w", a)
//	c := fmt.Errorf("third %w", b)
//	d := errors.First(c)
//	if a != d {
//		t.Fatal()
//	}
//	fmt.Println("")
//	x := Error[BadRequest]()
//
//	fmt.Println(x, x.T.Name)
//
//}

func TestExtract(t *testing.T) {
	a := errors.New("eloszki")
	b := fmt.Errorf("second %w", a)
	c := fmt.Errorf("third %w", b)

	for _, err := range errors.Extract(c) {
		fmt.Println(err)
	}
	//a1 := errors.Errorf("#1")
	//b1 := fmt.Errorf("#2 %w", a1)

	//x := errors.New("#4 %s %w %%nio %d", "elo", fmt.Errorf("#3 %w", fmt.Errorf("#2 %w", fmt.Errorf("#1"))), 13)
	//
	//y := errors.Trace(x)
	//fmt.Println("----")
	//fmt.Println(errors.First(x).String())
	//for _, err := range y {
	//	fmt.Println(err)
	//}

	//fmt.Println(errors.Unwrap(d1))
}

type Generic struct{}

func Foo() error {
	return Bar()
}

func Bar() error {
	err := errors.New("oh")
	fmt.Println(err.Message())
	return err
}
