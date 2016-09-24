package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	ttemplate "text/template"
)

var (
	cache map[string]*ttemplate.Template
	tmpl  *ttemplate.Template
	err   error
)

func init() {
	cache = make(map[string]*ttemplate.Template)
}

// New creates a template and caches it
func New(name string, template string) error {
	if _, ok := cache[name]; !ok {
		funcMap := ttemplate.FuncMap{
			"json": func(v interface{}) string {
				a, _ := json.Marshal(v)
				return string(a)
			},
		}
		// Create template, add Func map
		tmpl, err = ttemplate.New(name).Funcs(funcMap).Parse(template)
		if err != nil {
			return err
		}
		cache[name] = tmpl
	}
	return nil
}

// Exec returns a fully executed template substituted with a string map of row
func Exec(name string, row map[string]string) (string, error) {
	if _, ok := cache[name]; ok {
		tmpl = cache[name]
	} else {
		return "", fmt.Errorf("Exec called for template '%s' but not found in cache", name)
	}
	buf := bytes.NewBufferString("")
	err := tmpl.Execute(buf, row)

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}