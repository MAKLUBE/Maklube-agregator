package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
)

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	base := filepath.Join(dir, "base.tmpl")
	nav := filepath.Join(dir, "partials", "nav.tmpl")
	pagesDir := filepath.Join(dir, "pages")

	err := filepath.WalkDir(pagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		name := filepath.Base(path)
		ts, err := template.ParseFiles(base, nav, path)
		if err != nil {
			return fmt.Errorf("parse template %s: %w", name, err)
		}
		cache[name] = ts
		return nil
	})
	if err != nil {
		return nil, err
	}

	return cache, nil
}
