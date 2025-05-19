package report

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

// Finding holds data about a single fuzz trial.
type Finding struct {
	Service   string    `json:"service"`
	Method    string    `json:"method"`
	Payload   string    `json:"payload"`   // JSON of the mutated message
	Error     string    `json:"error"`     // gRPC error text (empty if OK)
	Timestamp time.Time `json:"timestamp"` // when the call was made
}

// WriteJSON writes findings to path as pretty-printed JSON.
func WriteJSON(findings []Finding, path string) error {
	data, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// WriteHTML writes findings into an HTML dashboard using tmplPath and writes to outPath.
func WriteHTML(findings []Finding, tmplPath, outPath string) error {
	// Use the base filename of tmplPath as the template name
	filename := filepath.Base(tmplPath)

	// Create a new template and register functions
	tmpl := template.New(filename).Funcs(template.FuncMap{
		"toJSON": func(v interface{}) (template.JS, error) {
			b, err := json.Marshal(v)
			return template.JS(b), err
		},
	})

	// Parse the template file
	tmpl, err := tmpl.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	// Create the output file
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute the named template and write to file
	return tmpl.ExecuteTemplate(f, filename, findings)
}
