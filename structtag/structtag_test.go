package structtag

import "testing"

func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Path", "PATH"},
		{"Driver", "DRIVER"},
		{"EnableBasic", "ENABLE_BASIC"},
		{"EnableOauth", "ENABLE_OAUTH"},
		{"AllowedOrigins", "ALLOWED_ORIGINS"},
		{"HTTPOnly", "HTTP_ONLY"},
		{"SameSite", "SAME_SITE"},
		{"MaxAgeSeconds", "MAX_AGE_SECONDS"},
		{"ClientID", "CLIENT_ID"},
		{"PublicBaseURL", "PUBLIC_BASE_URL"},
		{"BindAddress", "BIND_ADDRESS"},
		{"MaxSimpleKeys", "MAX_SIMPLE_KEYS"},
		{"ID", "ID"},
		{"A", "A"},
		{"", ""},
	}
	for _, c := range cases {
		if got := toSnakeCase(c.in); got != c.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
