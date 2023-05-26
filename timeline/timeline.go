package timeline

import (
	"os"
	"time"
)

type Moment struct {
    Instant time.Time
    FileName string

    format string
}

func NewMoment(format string) Moment {
}

func (m Moment) OpenOrCreate() (*os.File, error) {
}

func (m Moment) Next(n int) Moment {
}

func (m Moment) Back(n int) Moment {
}

func (m Moment) NextFunc(n int, f func(Moment) error) error {
}
