package single

import (
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
)

var envPattern = regexp.MustCompile(`(?P<KEY>[^=]+)=(?P<VALUJE>.*)`)

// secretKeyPattern matches env var names that commonly hold credentials;
// their values are masked in the UI to avoid leaking secrets on screen.
var secretKeyPattern = regexp.MustCompile(`(?i)(secret|passwd|password|token|api_?key|access_?key|private_?key|credential)`)

const redactedValue = "[REDACTED]"

type Env struct {
	*ui.Table
	data map[string]string
}

func NewEnv() *Env {
	p := ui.NewTable()
	p.Height = 4
	p.Width = colWidth[0]
	p.FgColor = ui.ThemeAttr("par.text.fg")
	p.Separator = false
	i := &Env{p, make(map[string]string)}
	i.BorderLabel = "Env"
	return i
}

func (w *Env) Set(allEnvs string) {
	envs := strings.Split(allEnvs, ";")
	w.Rows = [][]string{}
	for _, env := range envs {
		match := envPattern.FindStringSubmatch(env)
		if len(match) == 3 {
			key := match[1]
			value := match[2]
			if value != "" && secretKeyPattern.MatchString(key) {
				value = redactedValue
			}
			w.data[key] = value
			w.Rows = append(w.Rows, mkInfoRows(key, value)...)
		}
	}

	w.Height = len(w.Rows) + 2
}
