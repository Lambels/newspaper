package main

import (
	"context"
	"flag"
	"fmt"
	"os"
    "path/filepath"

	"github.com/Lambels/newspaper"
	"github.com/Lambels/newspaper/commands"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	fs := flag.NewFlagSet("np", flag.ExitOnError)
	flagReg, next, err := parseConfig(fs, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if next == nil {
		return
	}

	root := commands.NewRoot(fs)
	todayCmd := commands.NewToday(flagReg)
    findCmd := commands.NewFind(flagReg)

	root.Subcommands = []*ffcli.Command{
		todayCmd,
        findCmd,
	}

	if err := root.ParseAndRun(context.Background(), next); err != nil {
		fmt.Println(err)
	}
}

func parseConfig(fs *flag.FlagSet, args []string) (*commands.FlagRegister, []string, error) {
	flagReg := new(commands.FlagRegister)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}
    path := filepath.Join(homeDir, newspaper.DefaultConfigPath)
	cfgPath := fs.String("config", path, "set the path to the config file")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	cfg, err := newspaper.LoadConfig(*cfgPath)
	if err != nil {
		return nil, nil, err
	}

	flagReg.Config = cfg
	if fs.NFlag() != 0 && len(fs.Args()) == 0 {
		return flagReg, nil, nil
	}
	return flagReg, fs.Args(), nil
}
