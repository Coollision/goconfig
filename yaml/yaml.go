package yaml

import (
	"os"
	"path/filepath"

	"github.com/coollision/goconfig"
	"gopkg.in/yaml.v3"
)

func init() {
	f := goconfig.Fileformat{
		Extension:   ".yaml",
		Load:        LoadYAML,
		PrepareHelp: PrepareHelp,
	}
	goconfig.Formats = append(goconfig.Formats, f)
	f.Extension = ".yml"
	goconfig.Formats = append(goconfig.Formats, f)
}

// LoadYAML config file
func LoadYAML(config interface{}) (err error) {
	configFile := filepath.Join(goconfig.Path, goconfig.File)
	file, err := os.ReadFile(configFile)
	if os.IsNotExist(err) && !goconfig.FileRequired {
		err = nil
		return
	} else if err != nil {
		return
	}

	err = yaml.Unmarshal(file, config)

	if err != nil {
		return
	}

	return
}

// PrepareHelp return help string for this file format.
func PrepareHelp(config interface{}) (help string, err error) {
	var helpAux []byte
	helpAux, err = yaml.Marshal(&config)
	if err != nil {
		return
	}
	help = string(helpAux)
	return
}
