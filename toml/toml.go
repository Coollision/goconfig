// Package toml registers the TOML config file format for goconfig.
package toml

import (
	"os"
	"path/filepath"

	"github.com/coollision/goconfig"
	"github.com/pelletier/go-toml/v2"
)

func init() {
	f := goconfig.Fileformat{
		Extension:   ".toml",
		Load:        LoadTOML,
		PrepareHelp: PrepareHelp,
	}
	goconfig.Formats = append(goconfig.Formats, f)
}

// LoadTOML config file
func LoadTOML(config interface{}) (err error) {
	configFile := filepath.Join(goconfig.Path, goconfig.File)
	_, err = os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) && !goconfig.FileRequired {
			err = nil
		}
		return
	}
	var data []byte
	data, err = os.ReadFile(configFile)
	if err != nil {
		return
	}
	err = toml.Unmarshal(data, config)
	return
}

// PrepareHelp return help string for this file format.
func PrepareHelp(config interface{}) (help string, err error) {
	var byt []byte
	byt, err = toml.Marshal(config)
	if err != nil {
		return
	}
	help = string(byt)
	return
}
