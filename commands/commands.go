package commands

import (
	"flag"
	"fmt"

	"github.com/Lambels/newspaper"
)

type FlagRegister struct {
	Config *newspaper.Config
}

type GeneralContext struct {
	Verbose *bool
	Config  *newspaper.Config
}

func (c *FlagRegister) RegisterGlobal(fs *flag.FlagSet) GeneralContext {
	var ctx GeneralContext

	ctx.Config = c.Config
	ctx.Verbose = fs.Bool("v", false, "log with verbose output")

	return ctx
}

type TimeContext struct {
	Offset *int
}

func (c *FlagRegister) RegisterTime(fs *flag.FlagSet) TimeContext {
	var ctx TimeContext

	ctx.Offset = fs.Int("o", 0, "set action offset")

	return ctx
}

func LogVerbosef(verbose bool, format string, v ...any) {
	if verbose {
		fmt.Printf(format, v...)
	}
}
