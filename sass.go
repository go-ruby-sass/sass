// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2026, the go-ruby-sass/sass authors

// Package sass is a thin Ruby-facing adapter over the pure-Go Sass/SCSS
// compiler github.com/go-scss/scss. It exposes the modern sass-embedded gem
// surface — Sass.compile_string(source, syntax:, style:, load_paths:, ...) and
// Sass.compile(path, ...) returning {css, loaded_urls, source_map} — in a shape
// the rbgo binding wraps directly, plus the legacy SassC::Engine API (see the
// sassc subpackage).
//
// This module is intentionally minimal: it maps Ruby keyword arguments to the
// engine's typed options, calls the engine, and returns a result whose fields
// mirror the gem's CompileResult (`css`, `loaded_urls`, `source_map`).
package sass

import (
	engine "github.com/go-scss/scss"
)

// Syntax mirrors the gem's `syntax:` keyword ("scss", "indented", "css").
type Syntax string

const (
	SyntaxSCSS     Syntax = "scss"
	SyntaxIndented Syntax = "indented"
	SyntaxCSS      Syntax = "css"
)

// Style mirrors the gem's `style:` keyword ("expanded", "compressed").
type Style string

const (
	StyleExpanded   Style = "expanded"
	StyleCompressed Style = "compressed"
)

// Options mirrors the sass-embedded keyword arguments common to compile and
// compile_string. Zero values match the gem defaults (SCSS, expanded).
type Options struct {
	// Syntax is honored by CompileString; Compile infers it from the file
	// extension unless overridden.
	Syntax    Syntax
	Style     Style
	LoadPaths []string
	// Importer, when non-nil, resolves @use/@forward/@import URLs. It mirrors a
	// custom Ruby importer: given a URL it returns the source, a canonical URL,
	// and whether it handled the request.
	Importer func(url string) (source string, canonicalURL string, ok bool)
}

// CompileResult mirrors the gem's Sass::CompileResult: css, loaded_urls, and
// source_map (source maps are not yet emitted; see the engine residuals).
type CompileResult struct {
	CSS        string
	LoadedURLs []string
	SourceMap  string
}

func (o Options) engineOptions() *engine.Options {
	eo := &engine.Options{LoadPaths: o.LoadPaths}
	switch o.Syntax {
	case SyntaxIndented:
		eo.Syntax = engine.SyntaxIndented
	case SyntaxCSS:
		eo.Syntax = engine.SyntaxCSS
	default:
		eo.Syntax = engine.SyntaxSCSS
	}
	if o.Style == StyleCompressed {
		eo.Style = engine.Compressed
	}
	if o.Importer != nil {
		eo.Importer = engine.Importer(o.Importer)
	}
	return eo
}

// CompileString compiles Sass/SCSS source text, mirroring Sass.compile_string.
// This is the entry point Jekyll's jekyll-sass-converter calls.
func CompileString(source string, opts *Options) (*CompileResult, error) {
	if opts == nil {
		opts = &Options{}
	}
	res, err := engine.CompileString(source, opts.engineOptions())
	if err != nil {
		return nil, err
	}
	return &CompileResult{CSS: res.CSS, LoadedURLs: res.LoadedURLs, SourceMap: res.SourceMap}, nil
}

// Compile compiles a Sass/SCSS file, mirroring Sass.compile.
func Compile(path string, opts *Options) (*CompileResult, error) {
	if opts == nil {
		opts = &Options{}
	}
	eo := opts.engineOptions()
	// For Compile, syntax defaults to detection by extension; only force it when
	// the caller explicitly set a non-default syntax.
	if opts.Syntax == "" {
		eo.Syntax = engine.SyntaxSCSS
	}
	res, err := engine.Compile(path, eo)
	if err != nil {
		return nil, err
	}
	return &CompileResult{CSS: res.CSS, LoadedURLs: res.LoadedURLs, SourceMap: res.SourceMap}, nil
}
