package errors_test

import (
	"fmt"
	"testing"

	"github.com/sokool/errors"
)

func TestErrorf(t *testing.T) {
	type scenario struct {
		description                 string
		err                         *errors.Error // given
		name, code, message, string string        // expects
	}

	cases := []scenario{
		{
			description: "empty",
			err:         errors.Errorf(""),
			string:      "github.com/sokool/errors_test.TestErrorf#error_test.go@L20",
			name:        "github.com/sokool/errors_test.TestErrorf",
			code:        "error_test.go@L20",
		},
		{
			description: "just message",
			err:         errors.Errorf("hi there"),
			string:      "github.com/sokool/errors_test.TestErrorf#error_test.go@L27: hi there",
			name:        "github.com/sokool/errors_test.TestErrorf",
			code:        "error_test.go@L27",
			message:     "hi there",
		},
		{
			description: "message with arguments",
			err:         errors.Errorf("hi there %s", "man"),
			string:      "github.com/sokool/errors_test.TestErrorf#error_test.go@L35: hi there man",
			name:        "github.com/sokool/errors_test.TestErrorf",
			code:        "error_test.go@L35",
			message:     "hi there man",
		},
		{
			description: "just name",
			err:         errors.Errorf("test:"),
			string:      "test:",
			name:        "test",
		},
		{
			description: "just code",
			err:         errors.Errorf("#h6b7:"),
			string:      "#h6b7",
			code:        "h6b7",
		},
		{
			description: "with name,message",
			err:         errors.Errorf("test:hi there"),
			string:      "test: hi there",
			name:        "test",
			message:     "hi there",
		},
		{
			description: "with name,code,message",
			err:         errors.Errorf("email#h1:     invalid hostname      "),
			string:      "email#h1: invalid hostname",
			name:        "email",
			code:        "h1",
			message:     "invalid hostname",
		},
		{
			description: "with name,code,message from arguments",
			err:         errors.Errorf("%s#%s: invalid %s", "email", "h45", "username"),
			string:      "email#h45: invalid username",
			name:        "email",
			code:        "h45",
			message:     "invalid username",
		},
		{
			description: "with code,message",
			err:         errors.Errorf("#e87:failed"),
			string:      "#e87: failed",
			code:        "e87",
			message:     "failed",
		},
		{
			description: "with code,name signs in wrong place",
			err:         errors.Errorf("failed due abc: #triggered"),
			name:        "github.com/sokool/errors_test.TestErrorf",
			string:      "github.com/sokool/errors_test.TestErrorf#error_test.go@L85: failed due abc: #triggered",
			code:        "error_test.go@L85",
			message:     "failed due abc: #triggered",
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if s := c.err.Name(); s != c.name {
				t.Fatalf("expected name `%s`, got `%s`", c.name, s)
			}
			if s := c.err.Code(); s != c.code {
				t.Fatalf("expected code `%s`, got `%s`", c.code, s)
			}
			if s := c.err.Message(); s != c.message {
				t.Fatalf("expected message `%s`, got `%s`", c.message, s)
			}
			if s := c.err.Error(); s != c.string {
				t.Fatalf("expected string `%s`, got `%s`", c.string, s)
			}
		})
	}
}

func TestGetErr(t *testing.T) {
	a := errors.Errorf("eloszki")
	b := fmt.Errorf("second %w", a)
	c := fmt.Errorf("third %w", b)
	d := errors.Err(c)

	if a != d {
		t.Fatal()
	}

	//a1 := errors.Errorf("#1")
	//b1 := fmt.Errorf("#2 %w", a1)

	x := errors.Errorf("#4 %s %w %%nio %d", "elo", fmt.Errorf("#3 %w", fmt.Errorf("#2 %w", fmt.Errorf("#1"))), 13)

	y := errors.ErrL(x)
	fmt.Println("----")

	for _, err := range y {
		fmt.Println(err)
	}

	//fmt.Println(errors.Unwrap(d1))

}
