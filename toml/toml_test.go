package toml_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coollision/goconfig"
	_ "github.com/coollision/goconfig/toml"
)

type tomlConfig struct {
	LogLevel string `cfgDefault:"INFO"`
	Api      struct {
		Auth struct {
			EnableBasic bool `cfgDefault:"false"`
		}
	}
}

func TestLoadTOML_FileValuesApply(t *testing.T) {
	dir := t.TempDir()
	content := "LogLevel = \"trace\"\n\n[Api.Auth]\nEnableBasic = true\n"
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(content), 0o644); err != nil {
		t.Fatalf("write config.toml: %v", err)
	}

	goconfig.Path = dir
	goconfig.File = "config.toml"
	defer func() { goconfig.Path = "./"; goconfig.File = "" }()

	cfg := &tomlConfig{}
	if err := goconfig.Parse(cfg); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.LogLevel != "trace" {
		t.Errorf("LogLevel = %q, want trace", cfg.LogLevel)
	}
	if !cfg.Api.Auth.EnableBasic {
		t.Error("expected Api.Auth.EnableBasic to be true from config.toml")
	}
}

func TestLoadTOML_MissingFileFallsBackToDefaults(t *testing.T) {
	goconfig.Path = t.TempDir() // empty dir, no config.toml in it
	goconfig.File = "config.toml"
	defer func() { goconfig.Path = "./"; goconfig.File = "" }()

	cfg := &tomlConfig{}
	if err := goconfig.Parse(cfg); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.LogLevel != "INFO" {
		t.Errorf("LogLevel = %q, want INFO (cfgDefault)", cfg.LogLevel)
	}
}
