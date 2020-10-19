package fsm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		str        string
		want       string
		firstUpper bool
	}{
		{str: "some_str", want: "SomeStr", firstUpper: true},
		{str: "", want: "", firstUpper: true},
		{str: "F", want: "F", firstUpper: true},
		{str: "foo", want: "Foo", firstUpper: true},
		{str: "FooBar", want: "FooBar", firstUpper: true},
		{str: "fooBarBaz", want: "FooBarBaz", firstUpper: true},
		{str: "fooBar_baz", want: "FooBarBaz", firstUpper: true},
		{str: " foo_bar\n", want: "FooBar", firstUpper: true},
		{str: " foo-bar\t", want: "FooBar", firstUpper: true},
		{str: " foo bar\r", want: "FooBar", firstUpper: true},
		{str: "HTTP_status_code", want: "HttpStatusCode", firstUpper: true},

		{str: "some_str", want: "someStr", firstUpper: false},
		{str: "", want: "", firstUpper: false},
		{str: "F", want: "f", firstUpper: false},
		{str: "foo", want: "foo", firstUpper: false},
		{str: "FooBar", want: "fooBar", firstUpper: false},
		{str: "fooBarBaz", want: "fooBarBaz", firstUpper: false},
		{str: "fooBar_baz", want: "fooBarBaz", firstUpper: false},
		{str: " foo_bar\n", want: "fooBar", firstUpper: false},
		{str: " foo-bar\t", want: "fooBar", firstUpper: false},
		{str: " foo bar\r", want: "fooBar", firstUpper: false},
		{str: "HTTP_status_code", want: "httpStatusCode", firstUpper: false},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := ToCamelCase(tt.str, tt.firstUpper)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ToCamelCase() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
