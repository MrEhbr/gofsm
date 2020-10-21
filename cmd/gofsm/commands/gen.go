package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/MrEhbr/gofsm/fsm"
	"github.com/urfave/cli/v2"
)

var GenCommand = &cli.Command{
	Name:  "gen",
	Usage: "generates fsm",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "package",
			Aliases:     []string{"p"},
			Usage:       "package where struct is located",
			DefaultText: "default is current dir(.)",
			Value:       ".",
		},
		&cli.StringFlag{
			Name:     "struct",
			Aliases:  []string{"s"},
			Usage:    "struct name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "field",
			Aliases:  []string{"f"},
			Usage:    "state field of struct",
			Required: true,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "output file name",
			DefaultText: "default srcdir/<struct>_fsm.go",
		},
		&cli.StringFlag{
			Name:    "transitions",
			Aliases: []string{"t"},
			Usage:   "path to file with transitions",
		},
		&cli.BoolFlag{
			Name:    "noGenerate",
			Aliases: []string{"g"},
			Usage:   "don't put //go:generate instruction to the generated code",
		},
		&cli.StringFlag{
			Name:    "graph-output",
			Aliases: []string{"a"},
			Usage:   "path to transition graph file in dot format",
		},
	},
	Action: func(c *cli.Context) error {
		if c.IsSet("transitions") {
			path, err := filepath.Abs(c.String("transitions"))
			if err != nil {
				return err
			}

			if err := c.Set("transitions", path); err != nil {
				return cli.NewExitError(err, 1)

			}
		}

		if c.IsSet("graph-output") {
			path, err := filepath.Abs(c.String("graph-output"))
			if err != nil {
				return err
			}

			if err := c.Set("graph-output", path); err != nil {
				return cli.NewExitError(err, 1)

			}
		}

		path, err := filepath.Abs(c.String("package"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := c.Set("package", path); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := genAction(c); err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	},
}

func genAction(c *cli.Context) error {
	options := fsm.Options{
		InputPackage:          c.String("package"),
		Struct:                c.String("struct"),
		StateField:            c.String("field"),
		TransitionsFile:       c.String("transitions"),
		OutputFile:            c.String("output"),
		DisableGoGenerate:     c.Bool("noGenerate"),
		ActionGraphOutputFile: c.String("graph-output"),
	}

	if options.OutputFile == "" {
		options.OutputFile = fmt.Sprintf("%s_fsm.go", strings.ToLower(options.Struct))
	}

	if strings.ContainsRune(options.OutputFile, os.PathSeparator) {
		return fmt.Errorf("output file contains path separator")
	}

	if options.TransitionsFile != "" {
		path, err := filepath.Abs(options.TransitionsFile)
		if err != nil {
			return err
		}

		options.TransitionsFile = path
	}

	var outDir = filepath.Base(options.InputPackage)
	isDir, err := isDirectory(options.InputPackage)
	if err != nil {
		return err
	}

	if isDir {
		outDir = options.InputPackage
	}

	options.OutputFile = filepath.Join(outDir, options.OutputFile)

	g, err := fsm.NewGenerator(options)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := g.Generate(buf); err != nil {
		return err
	}

	if err := ioutil.WriteFile(options.OutputFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	if options.ActionGraphOutputFile != "" {
		buf := &bytes.Buffer{}
		if err := g.GenerateTransitionGraph(buf); err != nil {
			return err
		}

		if err := ioutil.WriteFile(options.ActionGraphOutputFile, buf.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
