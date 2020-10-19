package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/MrEhbr/gofsm/cmd/gofsm/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s - version %s\n", c.App.Name, version)
		fmt.Printf("  commit: \t%s\n", commit)
		fmt.Printf("  build date: \t%s\n", date)
		fmt.Printf("  build user: \t%s\n", builtBy)
		fmt.Printf("  go version: \t%s\n", runtime.Version())
	}

	app := &cli.App{
		Name:        "gofsm",
		Description: "gofsm generates finite state machine for struct",
		Usage:       "tool for generating fsm",
		Version:     version,
		Commands: []*cli.Command{
			commands.GenCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

// These values are private which ensures they can only be set with the build flags.
//nolint:unused
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)
