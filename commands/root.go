package commands

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewRoot() *ffcli.Command {
	return &ffcli.Command{
		Name:       "np",
		ShortUsage: "np [flags] <subcommand> [flags] [<arg>...]",
		Exec: func(context.Context, []string) error {
            return flag.ErrHelp
		},
	}
}
