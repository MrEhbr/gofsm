package fsm

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	type args struct {
		opt Options
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		inspect    func(t *testing.T, gen *Generator)
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "not package",
			args: func(t *testing.T) args {
				return args{
					Options{
						InputPackage: "/dev/null",
						Struct:       "Order",
						StateField:   "State",
					}}
			},

			wantErr: true,
		},
		{
			name: "not struct",
			args: func(t *testing.T) args {

				const src = `package order
				type Order int
				`
				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage: fname,
						Struct:       "Order",
						StateField:   "State",
					}}
			},

			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				if !strings.Contains(err.Error(), "is not struct") {
					t.Fatalf("wanted is not struct error, got %s", err)
				}
			},
		},
		{
			name: "invalid not state field",
			args: func(t *testing.T) args {
				const src = `package order
				type Order struct {
					ID int
				}
				`
				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage: fname,
						Struct:       "Order",
						StateField:   "State",
					}}
			},

			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				if !strings.Contains(err.Error(), "state type") {
					t.Fatalf("wanted state type is empty error, got %s", err)
				}
			},
		},
		{
			name: "state field type",
			args: func(t *testing.T) args {
				const src = `package fsm
				type Order struct {
					ID int
					State int
				}
				`

				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage: fname,
						Struct:       "Order",
						StateField:   "State",
					}}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				if !strings.Contains(err.Error(), "not enum") {
					t.Fatalf("wanted not enum error, got: %s", err)
				}
			},
		},
		{
			name: "transitions file exists",
			args: func(t *testing.T) args {
				const src = `package fsm
type StateType int

const (
	Created StateType = iota
	Started
	Finished
	Failed
)
type Order struct {
	ID int
	State StateType
}`

				const transitions = `[{"from": ["CREATED"],"to": "STARTED","event": "place_order"}]`
				f, err := ioutil.TempFile(t.TempDir(), "trs")
				if err != nil {
					t.Fatal(err)
				}
				f.WriteString(transitions)

				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage:    fname,
						Struct:          "Order",
						StateField:      "State",
						TransitionsFile: f.Name(),
					}}
			},
			wantErr: false,
		},
		{
			name: "transitions not file exists",
			args: func(t *testing.T) args {
				const src = `package fsm
type StateType int

const (
	Created StateType = iota
	Started
	Finished
	Failed
)
type Order struct {
	ID int
	State StateType
}`

				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage:    fname,
						Struct:          "Order",
						StateField:      "State",
						TransitionsFile: "/dev/null/transitions.json",
					}}
			},
			wantErr: true,
		},
		{
			name: "no errors json transitions",
			args: func(t *testing.T) args {
				const src = `package fsm
type StateType int

const (
	Created StateType = iota
	Started
	Finished
	Failed
)
type Order struct {
	ID int
	State StateType
}`

				const transitions = `[{"from": ["CREATED"],"to": "STARTED","event": "place_order"}]`
				f, err := ioutil.TempFile(t.TempDir(), "*_trs.json")
				if err != nil {
					t.Fatal(err)
				}
				f.WriteString(transitions)

				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage:    fname,
						Struct:          "Order",
						StateField:      "State",
						TransitionsFile: f.Name(),
					}}
			},
			wantErr: false,
			inspect: func(t *testing.T, gen *Generator) {
				if err := gen.struc.Validate(); err != nil {
					t.Fatalf("exept valid struct with transitions, got: %s", err)
				}
			},
		},
		{
			name: "no errors yaml transitions from is string",
			args: func(t *testing.T) args {
				const src = `package fsm
type StateType int

const (
	Created StateType = iota
	Started
	Finished
	Failed
)
type Order struct {
	ID int
	State StateType
}`

				const transitions = `
- from: CREATED
  to: STARTED
  event: place_order`
				f, err := ioutil.TempFile(t.TempDir(), "*_trs.yml")
				if err != nil {
					t.Fatal(err)
				}
				f.WriteString(transitions)

				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage:    fname,
						Struct:          "Order",
						StateField:      "State",
						TransitionsFile: f.Name(),
					}}
			},
			wantErr: false,
			inspect: func(t *testing.T, gen *Generator) {
				if err := gen.struc.Validate(); err != nil {
					t.Fatalf("exept valid struct with transitions, got: %s", err)
				}
			},
		},
		{
			name: "no errors yaml transitions from is array",
			args: func(t *testing.T) args {
				const src = `package fsm
type StateType int

const (
	Created StateType = iota
	Started
	Finished
	Failed
)
type Order struct {
	ID int
	State StateType
}`

				const transitions = `
- from:
  - CREATED
  to: STARTED
  event: place_order`
				f, err := ioutil.TempFile(t.TempDir(), "*_trs.yml")
				if err != nil {
					t.Fatal(err)
				}
				f.WriteString(transitions)

				fname, err := createPackage("fsm", t.TempDir(), []byte(src))
				if err != nil {
					t.Fatal(err)
				}

				return args{
					Options{
						InputPackage:    fname,
						Struct:          "Order",
						StateField:      "State",
						TransitionsFile: f.Name(),
					}}
			},
			wantErr: false,
			inspect: func(t *testing.T, gen *Generator) {
				if err := gen.struc.Validate(); err != nil {
					t.Fatalf("exept valid struct with transitions, got: %s", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := NewGenerator(tArgs.opt)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NewGenerator error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspect != nil {
				tt.inspect(t, got1)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func createPackage(pkg, dir string, src []byte) (string, error) {
	mod := fmt.Sprintf(`module %s

go 1.15`, pkg)
	if err := ioutil.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0666); err != nil {
		return "", err
	}

	fpath := filepath.Join(dir, pkg+".go")
	if err := ioutil.WriteFile(fpath, src, 0666); err != nil {
		return "", err
	}

	return fpath, nil
}
