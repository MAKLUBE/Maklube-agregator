package app

import (
	"fmt"
	"html/template"
	"path/filepath"
)

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(
			filepath.Join(dir, "base.tmpl"),
			filepath.Join(dir, "partials", "nav.tmpl"),
			page,
		)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		cache[name] = ts
	}

	return cache, nil
}
