// Package redact masks secret-like values before they are rendered or
// logged. It is shared by the single-view Env widget and the debug
// container dump so both paths apply identical redaction.
package redact

import (
	"regexp"
	"strings"
)

// Redacted is the placeholder shown in place of a secret value.
const Redacted = "[REDACTED]"

// secretKeyPattern matches env var names that commonly hold credentials.
var secretKeyPattern = regexp.MustCompile(`(?i)(secret|passwd|password|token|api_?key|access_?key|private_?key|credential|auth|dsn|cert|cookie|jwt)`)

// credentialValuePattern matches connection-string values that embed
// credentials regardless of the key name, e.g. postgres://user:pass@host/db.
var credentialValuePattern = regexp.MustCompile(`://[^/@\s]+:[^@\s]+@`)

// Env returns the value of a single KEY=VALUE pair, masked if the key looks
// secret-like or the value embeds credentials.
func Env(key, value string) string {
	if value == "" {
		return value
	}
	if secretKeyPattern.MatchString(key) || credentialValuePattern.MatchString(value) {
		return Redacted
	}
	return value
}

// EnvList masks secret-like values in a semicolon-joined list of KEY=VALUE
// pairs, the format stored in container meta under "[ENV-VAR]".
func EnvList(envs string) string {
	if envs == "" {
		return envs
	}
	parts := strings.Split(envs, ";")
	for i, env := range parts {
		key, value, found := strings.Cut(env, "=")
		if !found {
			continue
		}
		parts[i] = key + "=" + Env(key, value)
	}
	return strings.Join(parts, ";")
}
