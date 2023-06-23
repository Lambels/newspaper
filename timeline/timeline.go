package timeline

import (
	"os"
	"path/filepath"
	"time"
)

// Moment represents a moment or day in a linear format opposed to the non linear format of the hierarchial directory. It allows for
// linear movement through the timeline, on demand creation and relative movement which is key for recalculating formulas.
type Moment struct {
	Instant time.Time

	format   string
	root     string
	relative int

	dirExists  bool
	fileExists bool
}

func NewMoment(format string, root string, instant time.Time) Moment {
	return Moment{
		Instant:    instant,
		format:     format,
		root:       root,
		relative:   0,
		dirExists:  false,
		fileExists: false,
	}
}

func (m Moment) copy() Moment {
	return Moment{
		format:   m.format,
		root:     m.root,
		relative: m.relative,
	}
}

func (m *Moment) Recenter(leftright ...*Moment) {}

func (m *Moment) Exists() (bool, error) {
	switch {
	case m.fileExists: // we already know that the file exists.
		return true, nil
	default:
		if _, err := os.Stat(m.String()); os.IsNotExist(err) {
			return false, nil
		} else if err != nil {
			return false, err
		} else {
			m.dirExists = true
			m.fileExists = true
			return true, nil
		}
	}
}

func (m *Moment) OpenOrCreate() (*os.File, error) {
	if !m.dirExists {
		if err := os.MkdirAll(filepath.Dir(m.String()), 0755); err != nil {
			return nil, err
		}
		m.dirExists = true
	}

	file, err := os.Create(m.String())
	if err == nil {
		m.fileExists = true
	}
	return file, err
}

func (m Moment) Next(n int) Moment {
	moment := m.copy()
	next, ok := sameMonth(m.Instant, n)
	if !ok {
		moment.dirExists = m.dirExists
	}
	moment.Instant = next
	moment.relative += n
	return moment
}

// TODO: might need to pass pointer + figure out dirExists optimisation.
func (m Moment) NextFunc(n int, f func(Moment) error) error {
	for i := 0; i < n; i++ {
		if err := f(m.Next(i)); err != nil {
			return err
		}
	}
	return nil
}

func (m Moment) String() string {
	return filepath.Join(m.root, m.Instant.Format(m.format))
}

func sameMonth(now time.Time, offset int) (time.Time, bool) {
	next := now.AddDate(0, 0, offset)
	if next.Month() != now.Month() {
		return next, false
	}
	return next, true
}
