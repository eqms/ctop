package redact

import "testing"

func TestEnv(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{"password key", "POSTGRES_PASSWORD", "hunter2", Redacted},
		{"token key", "GITHUB_TOKEN", "ghp_abc123", Redacted},
		{"api key with underscore", "MY_API_KEY", "abc", Redacted},
		{"auth key", "BASIC_AUTH", "user:pass", Redacted},
		{"jwt key", "JWT_SIGNING", "xyz", Redacted},
		{"connection string value", "DATABASE_URL", "postgres://user:pass@db:5432/app", Redacted},
		{"plain url without credentials", "ENDPOINT", "https://example.com/path", "https://example.com/path"},
		{"harmless var", "PATH", "/usr/local/bin:/usr/bin", "/usr/local/bin:/usr/bin"},
		{"harmless term", "TERM", "xterm-256color", "xterm-256color"},
		{"empty value stays empty", "SECRET_KEY", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Env(tt.key, tt.value); got != tt.want {
				t.Errorf("Env(%q, %q) = %q, want %q", tt.key, tt.value, got, tt.want)
			}
		})
	}
}

func TestEnvList(t *testing.T) {
	in := "PATH=/usr/bin;POSTGRES_PASSWORD=hunter2;DATABASE_URL=postgres://u:p@db/app;LANG=C.UTF-8;NOEQUALS"
	want := "PATH=/usr/bin;POSTGRES_PASSWORD=" + Redacted + ";DATABASE_URL=" + Redacted + ";LANG=C.UTF-8;NOEQUALS"
	if got := EnvList(in); got != want {
		t.Errorf("EnvList() = %q, want %q", got, want)
	}
	if got := EnvList(""); got != "" {
		t.Errorf("EnvList(\"\") = %q, want empty", got)
	}
}
