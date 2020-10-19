package fsm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/constant"
	"go/types"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

// Options for NewGenerator constructor
type Options struct {
	// InputPackage is an import path or a relative path of the package that contains the source struct
	InputPackage string
	// OutputFile name for output file
	OutputFile string
	// Struct is name of struct
	Struct string
	// StateField is name of struct field that indicate state
	StateField string
	// TransitionsFile is path to file which hold supported transitions
	TransitionsFile string
	// DisableGoGenerate don't put go generate
	DisableGoGenerate bool
}

// Generator generates finite state machine
type Generator struct {
	Options

	struc      Struct
	srcPackage *packages.Package

	fsmTemplate *template.Template
}

// NewGenerator returns a pointer to Generator
func NewGenerator(opt Options) (*Generator, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:  filepath.Dir(opt.InputPackage),
	}

	pkgs, err := packages.Load(cfg, opt.InputPackage)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(pkgs))
	}

	funcs := map[string]interface{}{
		"to_camel":  ToCamelCase,
		"join":      strings.Join,
		"rel":       filepath.Rel,
		"dir":       filepath.Dir,
		"base":      filepath.Base,
		"path_join": filepath.Join,
	}

	tpl, err := template.New("fsm").Funcs(funcs).Parse(fsmTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	gen := &Generator{
		Options:     opt,
		srcPackage:  pkgs[0],
		fsmTemplate: tpl,
	}

	if err := gen.fillStruct(); err != nil {
		return nil, err
	}

	if opt.TransitionsFile != "" {
		data, err := ioutil.ReadFile(opt.TransitionsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read transitions file: %w", err)
		}

		if err := json.Unmarshal(data, &gen.struc.Transitions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transitions file: %w", err)
		}
	}

	if err := gen.struc.Validate(); err != nil {
		return nil, err
	}

	return gen, nil
}

// fillStruct find info about struct and state values for generator
func (g *Generator) fillStruct() error {
	g.struc = Struct{
		Name:        g.Struct,
		StateField:  g.StateField,
		StateValues: make([]StateValue, 0, 10),
	}

	for _, def := range g.srcPackage.TypesInfo.Defs {
		typ, ok := def.(*types.TypeName)
		if !ok {
			continue
		}

		if def.Name() != g.Struct {
			continue
		}

		st, ok := typ.Type().Underlying().(*types.Struct)
		if !ok {
			return fmt.Errorf("%s is not struct", g.Struct)
		}

		for i := 0; i < st.NumFields(); i++ {
			if f := st.Field(i); f.Name() == g.StateField {
				typ, ok := f.Type().(*types.Named)
				if !ok {
					return fmt.Errorf("%s.%s is not enum", def.Name(), f.Name())
				}

				g.struc.StateType = typ.Obj().Name()
				break
			}
		}

		// we found needed struct so no need to loop anymore
		break
	}

	// find values for field type
	for _, def := range g.srcPackage.TypesInfo.Defs {
		val, ok := def.(*types.Const)
		if !ok {
			continue
		}

		typ, ok := val.Type().(*types.Named)
		if !ok {
			continue
		}

		if typ.Obj().Name() != g.struc.StateType {
			continue
		}

		if val.Name() == "_" {
			continue
		}

		valStr := val.Val().String()
		if val.Val().Kind() == constant.String {
			valStr, _ = strconv.Unquote(valStr)
		}

		g.struc.StateValues = append(g.struc.StateValues, StateValue{
			Name: val.Name(),
			Val:  valStr,
		})
	}

	return nil
}

// Generate generates code using template
func (g *Generator) Generate(w io.Writer) error {
	buf := bytes.NewBuffer([]byte{})

	if err := g.fsmTemplate.Execute(buf, map[string]interface{}{
		"Struct":      g.struc,
		"Transitions": g.struc.Transitions,
		"Package":     g.srcPackage,
		"Options":     g.Options,
	}); err != nil {
		return err
	}

	processedSource, err := imports.Process(g.Options.OutputFile, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	_, err = w.Write(processedSource)
	return err
}
