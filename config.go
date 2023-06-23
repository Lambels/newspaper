package newspaper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultConfigPath string = "/.config/newspaper.json"

type Config struct {
	Root          string `json:"root"`
	FileExtension string `json:"file_extension"`

	RemindersFirst bool `json:"reminders_first"`

	YearFromat  string `json:"year_format"`
	MonthFromat string `json:"month_format"`
	DayFormat   string `json:"day_format"`

	StorageFilePath string `json:"storage_file_path"`
}

func (c *Config) TimeFormat() string {
    return filepath.Join(c.YearFromat, c.MonthFromat, c.DayFormat)
}

func (c *Config) FileFormat() string {
	return fmt.Sprintf(
		"%v.%v",
        c.TimeFormat(),
		c.FileExtension,
	)
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}
	cfg.Root = filepath.Clean(cfg.Root)

	return cfg, nil
}
