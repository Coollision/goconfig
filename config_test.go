package goconfig

import "testing"

// Api and Auth mirror the real-world shape that exposed the bug: booleans
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

type Api struct {
	Cors struct {
		AllowedOrigins string `cfgDefault:""`
	}
	Auth Auth
}

type TestConfig struct {
	Api Api
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

	if !cfg.Api.Auth.EnableBasic {
		t.Error("expected Api.Auth.EnableBasic to be true from API_AUTH_ENABLE_BASIC")
	}
	if cfg.Api.Auth.EnableOauth {
		t.Error("expected Api.Auth.EnableOauth to be false from API_AUTH_ENABLE_OAUTH (its cfgDefault is true)")
	}
	if cfg.Api.Cors.AllowedOrigins != "https://example.test" {
		t.Errorf("expected Api.Cors.AllowedOrigins to be set from API_CORS_ALLOWED_ORIGINS, got %q", cfg.Api.Cors.AllowedOrigins)
	}
	if !cfg.Api.Auth.Cookie.Secure {
		t.Error("expected Api.Auth.Cookie.Secure to be true from API_AUTH_COOKIE_SECURE")
	}
}

func TestParse_DefaultsHoldWhenEnvUnset(t *testing.T) {
	cfg := &TestConfig{}
	if err := Parse(cfg); err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if cfg.Api.Auth.EnableBasic {
		t.Error("expected Api.Auth.EnableBasic to default to false")
	}
	if !cfg.Api.Auth.EnableOauth {
		t.Error("expected Api.Auth.EnableOauth to default to true")
	}
	if cfg.Api.Cors.AllowedOrigins != "" {
		t.Errorf("expected Api.Cors.AllowedOrigins to default to empty, got %q", cfg.Api.Cors.AllowedOrigins)
	}
}
