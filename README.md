# goconfig
[![CI](https://github.com/coollision/goconfig/actions/workflows/ci.yml/badge.svg)](https://github.com/coollision/goconfig/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/coollision/goconfig)](https://goreportcard.com/report/github.com/coollision/goconfig)
[![Maintainability](https://api.codeclimate.com/v1/badges/f11c9124505888c4c8db/maintainability)](https://codeclimate.com/github/coollision/goconfig/maintainability)
[![Go Reference](https://pkg.go.dev/badge/github.com/coollision/goconfig.svg)](https://pkg.go.dev/github.com/coollision/goconfig)
[![Go project version](https://badge.fury.io/go/github.com%2Fcoollision%2Fgoconfig.svg)](https://badge.fury.io/go/github.com%2Fcoollision%2Fgoconfig)[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)


goconfig uses a struct as input and populates the fields of this struct with parameters from command line, environment variables and configuration file.

## Install

```
go get github.com/coollision/goconfig
```

## Example

```go
package main

import (
	"fmt"

	"github.com/coollision/goconfig"
)

/*
step 1: Declare your configuration struct,
it may or may not contain substructures.
*/

type mongoDB struct {
	Host string `cfgDefault:"example.com" cfgRequired:"true"`
	Port int    `cfgDefault:"999"`
}

type configTest struct {
	Domain    string
	DebugMode bool `json:"db" cfg:"db" cfgDefault:"false"`
	MongoDB   mongoDB
	IgnoreMe  string `cfg:"-"`
}

func main() {

	// step 2: Instantiate your structure.
	config := new(configTest)

	// step 3: Pass the instance pointer to the parser
	err := goconfig.Parse(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
	   The parser populated your struct with the data
	   it took from environment variables and command
	   line, now you can use it.
	*/

	println("config.Domain......:", config.Domain)
	println("config.DebugMode...:", config.DebugMode)
	println("config.MongoDB.Host:", config.MongoDB.Host)
	println("config.MongoDB.Port:", config.MongoDB.Port)
}
```

With the example above try environment variables like *$DOMAIN* or *$MONGO_DB_HOST* and run the example again to see what happens.

You can also try using parameters on the command line, try -h to see the help.

## Config files

Environment variables and flags work without any extra setup. To read a configuration file, register the file format you need by importing its package for side effects:

```go
import (
	"github.com/coollision/goconfig"
	_ "github.com/coollision/goconfig/yaml" // .yaml and .yml
	_ "github.com/coollision/goconfig/toml" // .toml
)
```

Only import the formats you use. The file name is set with `goconfig.File` or the `GO_CONFIG_FILE` environment variable, and the directory with `goconfig.Path` (default `./`) or `GO_CONFIG_PATH`. The format is chosen from the file extension; a file with an unregistered extension makes `Parse` return `goconfig.ErrFileFormatNotDefined`. A missing file is ignored unless `goconfig.FileRequired` is `true`.

Values are applied in this order, later sources winning: `cfgDefault`, then the config file, then environment variables, then command-line flags.

`ParseAndWatch` loads the config file and validates it, and it does not read environment variables or flags. Set `goconfig.WatchConfigFile = true` before calling it to reload the file whenever it changes. Each reload sends a Unix timestamp on the returned channel, and watch errors are sent on the error channel. Use `Parse` if you need environment variables and flags.

## Prefixes

`goconfig.PrefixEnv` and `goconfig.PrefixFlag` add a prefix to every generated name.

- With `goconfig.PrefixEnv = "MYAPP"`, the field `Domain` is read from `MYAPP_DOMAIN`, nested fields are joined with underscores (`MYAPP_MONGO_DB_HOST`), and the file variables become `MYAPP_GO_CONFIG_FILE` and `MYAPP_GO_CONFIG_PATH`.
- With `goconfig.PrefixFlag = "app"`, the same field becomes the flag `-app_domain`.

## Struct tags

| Tag | Purpose |
| --- | --- |
| `cfg:"name"` | Name used for the environment variable and flag, e.g. `cfg:"db"` gives `$DB` and `-db`. Without it, the Go field name is converted to upper snake case (`MongoDB` becomes `MONGO_DB`, `EnableBasic` becomes `ENABLE_BASIC`). |
| `cfg:"-"` | Ignores the field completely. |
| `cfgDefault:"value"` | Value used when nothing else sets the field. For booleans, `"true"` or `"t"` means true. |
| `cfgHelper:"text"` | Usage text shown for the command-line flag in `-h` output. |
| `cfgRequired:"true"` | `Parse` returns an error if a string, int or float field is left empty or zero. |

Unexported fields are skipped. Nested structs are supported and their names are joined with underscores.

## Contributing

- Fork the repo on GitHub
- Clone the project to your own machine
- Create a *branch* with your modifications `git checkout -b fantastic-feature`.
- Then _commit_ your changes `git commit -m 'Implementation of new fantastic feature'`
- Make a _push_ to your _branch_ `git push origin fantastic-feature`.
- Submit a **Pull Request** so that we can review your changes

Please use [conventional commit](https://www.conventionalcommits.org/) messages and PR titles (`feat:`, `fix:`, `chore:`, ...). Releases and the changelog are generated by [release-please](https://github.com/googleapis/release-please) from these, and PRs are squash-merged, so the PR title becomes the commit message on `master`.
