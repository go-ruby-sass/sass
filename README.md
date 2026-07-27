# go-ruby-sass/sass

[![CI](https://github.com/go-ruby-sass/sass/actions/workflows/ci.yml/badge.svg)](https://github.com/go-ruby-sass/sass/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-ruby-sass/sass.svg)](https://pkg.go.dev/github.com/go-ruby-sass/sass)
[![License: BSD-3-Clause](https://img.shields.io/badge/License-BSD--3--Clause-blue.svg)](LICENSE)

A **thin Ruby-facing adapter** over the pure-Go, CGO-free Sass/SCSS compiler
**[github.com/go-scss/scss](https://github.com/go-scss/scss)**. It exposes the
modern **`sass-embedded`** gem surface (and the legacy **`sassc`** surface) in a
shape the [rbgo](https://github.com/go-embedded-ruby) binding wraps directly, so
pure-Go Ruby can compile Sass with no cgo and no external `sass` binary.

All language work and the dart-sass differential gate live in **go-scss/scss**;
this module only maps Ruby keyword arguments to the engine's typed options and
returns a result mirroring the gem's `Sass::CompileResult`.

## Gem API surface

### sass-embedded (`Sass`)

```go
import sass "github.com/go-ruby-sass/sass"

// Sass.compile_string(source, syntax:, style:, load_paths:, importer:)
res, err := sass.CompileString(src, &sass.Options{
    Syntax:    sass.SyntaxSCSS,       // scss | indented | css
    Style:     sass.StyleCompressed,  // expanded | compressed
    LoadPaths: []string{"_sass"},
})
_ = res.CSS          // => css
_ = res.LoadedURLs   // => loaded_urls
_ = res.SourceMap    // => source_map (empty; see engine residuals)

// Sass.compile(path, ...)
res, err = sass.Compile("styles/main.scss", nil)
```

`CompileString` is the exact entry point **Jekyll's `jekyll-sass-converter`**
calls — that path is covered by a dedicated test.

### sassc (`SassC::Engine`)

```go
import "github.com/go-ruby-sass/sass/sassc"

// SassC::Engine.new(template, style: :compressed).render
css, err := sassc.NewEngine(template, sassc.EngineOptions{
    Style: sassc.StyleCompressed,
}).Render()
```

## Ruby mapping (for the rbgo binding)

| Ruby | Go |
| --- | --- |
| `Sass.compile_string(src, **kw)` | `sass.CompileString(src, opts)` |
| `Sass.compile(path, **kw)` | `sass.Compile(path, opts)` |
| result `.css` / `.loaded_urls` / `.source_map` | `CompileResult.CSS` / `.LoadedURLs` / `.SourceMap` |
| `syntax:` / `style:` / `load_paths:` | `Options.Syntax` / `.Style` / `.LoadPaths` |
| custom importer | `Options.Importer func(url) (src, canonical, ok)` |
| `SassC::Engine.new(t, style:).render` | `sassc.NewEngine(t, opts).Render()` |

## Residuals

- `source_map` is not yet emitted (the engine does not produce source maps yet).
- `sassc` `:nested`/`:compact` styles fall back to `expanded` (the engine emits
  `expanded` and `compressed`).
- Output fidelity is exactly that of the underlying engine — see the
  [go-scss/scss residuals](https://github.com/go-scss/scss#honest-residuals)
  (notably CSS Color 4 percentage serialization for some computed colors).

## License

BSD-3-Clause. Copyright (c) 2026, the go-ruby-sass/sass authors.
