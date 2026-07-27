// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2026, the go-ruby-sass/sass authors

// Package sassc mirrors the legacy SassC gem surface — SassC::Engine.new(
// template, style:).render — layered over the same pure-Go engine, since real
// applications still use both the modern sass-embedded and the older sassc API.
package sassc

import (
	sass "github.com/go-ruby-sass/sass"
)

// Style mirrors SassC's `:style` option (:expanded, :nested, :compact,
// :compressed). The pure-Go engine implements expanded and compressed; :nested
// and :compact fall back to expanded (their legacy layouts are not emitted).
type Style string

const (
	StyleExpanded   Style = "expanded"
	StyleNested     Style = "nested"
	StyleCompact    Style = "compact"
	StyleCompressed Style = "compressed"
)

// EngineOptions mirrors SassC::Engine.new keyword options.
type EngineOptions struct {
	Style     Style
	Syntax    string // "scss" (default) or "sass"/"indented"
	LoadPaths []string
}

// Engine mirrors SassC::Engine.
type Engine struct {
	template string
	opts     EngineOptions
}

// NewEngine mirrors SassC::Engine.new(template, **options).
func NewEngine(template string, opts EngineOptions) *Engine {
	return &Engine{template: template, opts: opts}
}

// Render mirrors SassC::Engine#render, returning the compiled CSS.
func (e *Engine) Render() (string, error) {
	o := &sass.Options{LoadPaths: e.opts.LoadPaths}
	if e.opts.Style == StyleCompressed {
		o.Style = sass.StyleCompressed
	}
	switch e.opts.Syntax {
	case "sass", "indented":
		o.Syntax = sass.SyntaxIndented
	}
	res, err := sass.CompileString(e.template, o)
	if err != nil {
		return "", err
	}
	return res.CSS, nil
}
