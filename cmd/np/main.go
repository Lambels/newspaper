package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Lambels/newspaper"
	"github.com/Lambels/newspaper/commands"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
    flagReg, next, err := parseConfig(os.Args[1:])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if next == nil {
        return
    }

    root := commands.NewRoot()
    todayCmd := commands.NewToday(flagReg)

    root.Subcommands = []*ffcli.Command{
        todayCmd,
    }

    root.ParseAndRun(context.Background(), next)
}

func parseConfig(args []string) (*commands.FlagRegister, []string, error) {
	flagReg := new(commands.FlagRegister)

	fs := flag.NewFlagSet("np", flag.ExitOnError)
	ctx := flagReg.RegisterGlobal(fs)
	cfgPath := fs.String("config", "~/.config/newspaper.json", "set the path to the config file")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	cfg, err := newspaper.LoadConfig(*cfgPath)
	var pathErr os.PathError
	switch {
	case errors.As(err, &pathErr):
		commands.LogVerbosef(*ctx.Verbose, "could not open config file at: %s , using default config", pathErr.Path)
	case err != nil:
		return nil, nil, err
	}
    
    flagReg.Config = cfg
    if fs.NFlag() == 0 && len(fs.Args()) == 0 {
        return flagReg, nil, nil    
    }
    return flagReg, fs.Args(), nil
}
