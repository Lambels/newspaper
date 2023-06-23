package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Lambels/newspaper"
	"github.com/Lambels/newspaper/timeline"
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
	Root   *timeline.Moment
}

func (c *FlagRegister) RegisterTime(fs *flag.FlagSet) TimeContext {
	var ctx TimeContext

	moment := timeline.NewMoment(c.Config.FileFormat(), c.Config.Root, time.Now())
	fs.Func("r", "set the root of the action", func(s string) error {
		time, err := time.Parse(c.Config.TimeFormat(), s)
		if err != nil {
			return err
		}
		moment.Instant = time
		return nil
	})
	ctx.Offset = fs.Int("o", 0, "set action offset")
	ctx.Root = &moment

	return ctx
}

func logVerbosef(verbose bool, format string, v ...any) {
	if verbose {
		fmt.Printf(format, v...)
	}
}

func createOnForce(force bool, moment *timeline.Moment) (*os.File, bool, error) {
	if !force {
		return nil, false, nil
	}

	if exists, err := moment.Exists(); err == nil && !exists {
		f, err := moment.OpenOrCreate()
		return f, exists, err
	} else {
		return nil, false, err
	}
}
