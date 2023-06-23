package commands

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/Lambels/newspaper/timeline"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewToday(flagReg *FlagRegister) *ffcli.Command {
	fs := flag.NewFlagSet("np today", flag.ExitOnError)
	force := fs.Bool("f", false, "optionally creates the file")
	ctx := flagReg.RegisterGlobal(fs)

	return &ffcli.Command{
		Name:       "today",
		ShortUsage: "today",
		ShortHelp:  "outputs path to todays file or optionally creates it",
		FlagSet:    fs,
		Exec: func(context.Context, []string) error {
			moment := timeline.NewMoment(ctx.Config.FileFormat(), ctx.Config.Root, time.Now())
			fmt.Println(moment.String())

			_, existed, err := createOnForce(*force, &moment)
			if err == nil && existed {
				logVerbosef(*ctx.Verbose, "created todays file\n")
			}
			return err
		},
	}
}
