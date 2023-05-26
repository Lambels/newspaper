package newspaper

import (
	"encoding/json"
	"os"
)

const defaultConfigPath string = ""

var defaultConfig Config = Config{
    root:           "",
    fileExtension:  "md",
    remindersFirst: true,
    yearFromat:     "YYYY",
    monthFromat:    "MM",
    dayFormat:      "DD",
}

type Config struct {
	root          string
	fileExtension string

	remindersFirst bool

	yearFromat  string
	monthFromat string
	dayFormat   string

    storageFilePath string
}

func LoadConfig(path string) (*Config, error) {
    if path == "" {
        path = defaultConfigPath
    }

    file, err := os.Open(path)
    if err != nil {
        return &defaultConfig, err
    }

    cfg := new(Config)
    if err := json.NewDecoder(file).Decode(cfg); err != nil {
        return &defaultConfig, err
    }

    return cfg, nil
}
