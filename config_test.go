package goconfig

import "testing"

// API and Auth mirror the real-world shape that exposed the bug: booleans
// and strings nested two levels deep, alongside sibling fields at the same
// depth whose (single-word) names happened to already work. Before
// structtag's camelCase splitting fix, EnableBasic/EnableOauth/
// AllowedOrigins below silently ignored their environment variables while
// Cookie.Secure (single-word leaf) picked its up fine — see
// structtag.toSnakeCase's doc comment for the full story.
type Auth struct {
	EnableBasic bool `cfgDefault:"false"`
	EnableOauth bool `cfgDefault:"true"`
	Cookie      struct {
		Secure bool `cfgDefault:"false"`
	}
}

type API struct {
	Cors struct {
		AllowedOrigins string `cfgDefault:""`
	}
	Auth Auth
}

type TestConfig struct {
	API API
}

func TestParse_MultiWordNestedFieldsBindFromEnv(t *testing.T) {
	t.Setenv("API_AUTH_ENABLE_BASIC", "true")
	t.Setenv("API_AUTH_ENABLE_OAUTH", "false")
	t.Setenv("API_CORS_ALLOWED_ORIGINS", "https://example.test")
	t.Setenv("API_AUTH_COOKIE_SECURE", "true")

	cfg := &TestConfig{}
	if err := Parse(cfg); err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if !cfg.API.Auth.EnableBasic {
		t.Error("expected API.Auth.EnableBasic to be true from API_AUTH_ENABLE_BASIC")
	}
	if cfg.API.Auth.EnableOauth {
		t.Error("expected API.Auth.EnableOauth to be false from API_AUTH_ENABLE_OAUTH (its cfgDefault is true)")
	}
	if cfg.API.Cors.AllowedOrigins != "https://example.test" {
		t.Errorf("expected API.Cors.AllowedOrigins to be set from API_CORS_ALLOWED_ORIGINS, got %q", cfg.API.Cors.AllowedOrigins)
	}
	if !cfg.API.Auth.Cookie.Secure {
		t.Error("expected API.Auth.Cookie.Secure to be true from API_AUTH_COOKIE_SECURE")
	}
}

func TestParse_DefaultsHoldWhenEnvUnset(t *testing.T) {
	cfg := &TestConfig{}
	if err := Parse(cfg); err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if cfg.API.Auth.EnableBasic {
		t.Error("expected API.Auth.EnableBasic to default to false")
	}
	if !cfg.API.Auth.EnableOauth {
		t.Error("expected API.Auth.EnableOauth to default to true")
	}
	if cfg.API.Cors.AllowedOrigins != "" {
		t.Errorf("expected API.Cors.AllowedOrigins to default to empty, got %q", cfg.API.Cors.AllowedOrigins)
	}
}
