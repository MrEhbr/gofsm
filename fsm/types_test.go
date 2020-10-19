package fsm

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStruct_FindValue(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		init func(t *testing.T) Struct

		args func(t *testing.T) args

		want1 string
	}{
		{
			name: "not found",
			init: func(*testing.T) Struct {
				return Struct{
					StateValues: []StateValue{{Name: "Foo", Val: "foo"}},
				}
			},
			args: func(*testing.T) args {
				return args{"bar"}
			},

			want1: "",
		},
		{
			name: "val in snake case snake ",
			init: func(*testing.T) Struct {
				return Struct{
					StateValues: []StateValue{{Name: "FooBar", Val: "foo_bar"}},
				}
			},
			args: func(*testing.T) args {
				return args{"foo_bar"}
			},

			want1: "FooBar",
		},
		{
			name: "val in with spaces ",
			init: func(*testing.T) Struct {
				return Struct{
					StateValues: []StateValue{{Name: "FooBar", Val: "foo_bar"}},
				}
			},
			args: func(*testing.T) args {
				return args{"foo bar"}
			},

			want1: "FooBar",
		},
		{
			name: "find by value",
			init: func(*testing.T) Struct {
				return Struct{
					StateValues: []StateValue{{Name: "FooBar", Val: "baz"}},
				}
			},
			args: func(*testing.T) args {
				return args{"baz"}
			},

			want1: "FooBar",
		},
		{
			name: "only spaces",
			init: func(*testing.T) Struct {
				return Struct{
					StateValues: []StateValue{{Name: "FooBar", Val: "baz"}},
				}
			},
			args: func(*testing.T) args {
				return args{"      "}
			},

			want1: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1 := receiver.FindValue(tArgs.str)

			if diff := cmp.Diff(tt.want1, got1); diff != "" {
				t.Errorf("FindValue() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTransitions_Actions(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) Transitions
		inspect func(r Transitions, t *testing.T) //inspects receiver after test run

		want1 []string
	}{
		{
			name:  "empty actions",
			init:  func(*testing.T) Transitions { return Transitions{} },
			want1: nil,
		},
		{
			name: "all unique",
			init: func(*testing.T) Transitions {
				return Transitions{
					{Actions: []string{"foo"}, BeforeActions: []string{"bar"}},
					{Actions: []string{"baz"}, BeforeActions: []string{"foobar"}},
				}
			},
			want1: []string{"foo", "bar", "baz", "foobar"},
		},
		{
			name: "have duplicates",
			init: func(*testing.T) Transitions {
				return Transitions{
					{Actions: []string{"foo"}, BeforeActions: []string{"bar"}},
					{Actions: []string{"foo"}, BeforeActions: []string{"bar"}},
				}
			},
			want1: []string{"foo", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.Actions()

			// sort slices, because for tests wee care about values
			sort.Strings(got1)
			sort.Strings(tt.want1)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if diff := cmp.Diff(tt.want1, got1); diff != "" {
				t.Errorf("Actions() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
