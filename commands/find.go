package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewFind(flagReg *FlagRegister) *ffcli.Command {
    fs := flag.NewFlagSet("np find", flag.ExitOnError)
    force := fs.Bool("f", false, "force the creation of the todo you want to find")
    ctx := flagReg.RegisterGlobal(fs)
    tctx := flagReg.RegisterTime(fs)

    return &ffcli.Command{
        Name: "find",
        ShortUsage: "",
        ShortHelp: "",
        LongHelp: "",
        FlagSet: fs,
        Exec: func(context.Context, []string) error {
            next := tctx.Root.Next(*tctx.Offset)
            fmt.Println(next.String())
            _, existed, err := createOnForce(*force, &next)
            if err == nil && existed {
				logVerbosef(*ctx.Verbose, "created file\n")
            }

            return err
        },
    }
}
