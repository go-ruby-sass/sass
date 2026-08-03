// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2026, the go-ruby-sass/sass authors

package sass_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	sass "github.com/go-ruby-sass/sass"
	"github.com/go-ruby-sass/sass/sassc"
)

func TestCompileStringDefault(t *testing.T) {
	res, err := sass.CompileString(".a { b: 1 + 1; }", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.CSS != ".a {\n  b: 2;\n}\n" {
		t.Errorf("got %q", res.CSS)
	}
}

func TestCompileStringCompressed(t *testing.T) {
	res, err := sass.CompileString(".a{b:1}", &sass.Options{Style: sass.StyleCompressed})
	if err != nil {
		t.Fatal(err)
	}
	if res.CSS != ".a{b:1}\n" {
		t.Errorf("got %q", res.CSS)
	}
}

func TestCompileStringIndented(t *testing.T) {
	res, err := sass.CompileString(".a\n  b: 1\n", &sass.Options{Syntax: sass.SyntaxIndented})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, "b: 1") {
		t.Errorf("got %q", res.CSS)
	}
}

func TestCompileStringError(t *testing.T) {
	if _, err := sass.CompileString(".a{b:$undef}", nil); err == nil {
		t.Error("expected error")
	}
}

func TestCompileStringLoadedURLs(t *testing.T) {
	imp := func(url string) (string, string, bool) {
		if url == "part" {
			return "$c: #123;", "part", true
		}
		return "", "", false
	}
	res, err := sass.CompileString("@use \"part\"; .a{c: part.$c}", &sass.Options{Importer: imp})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, "#123") {
		t.Errorf("got %q", res.CSS)
	}
	if len(res.LoadedURLs) == 0 {
		t.Error("expected loaded_urls")
	}
}

func TestCompileFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "in.scss")
	if err := os.WriteFile(p, []byte(".x { y: 2 * 3px; }"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := sass.Compile(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, "y: 6px") {
		t.Errorf("got %q", res.CSS)
	}
}

func TestCompileFileMissing(t *testing.T) {
	if _, err := sass.Compile("nope.scss", nil); err == nil {
		t.Error("expected error")
	}
}

func TestCSSSyntaxOption(t *testing.T) {
	res, err := sass.CompileString(".a{b:1}", &sass.Options{Syntax: sass.SyntaxCSS})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, "b: 1") {
		t.Errorf("got %q", res.CSS)
	}
}

func TestSassCEngine(t *testing.T) {
	css, err := sassc.NewEngine(".a { b: 1 + 1; }", sassc.EngineOptions{}).Render()
	if err != nil {
		t.Fatal(err)
	}
	if css != ".a {\n  b: 2;\n}\n" {
		t.Errorf("got %q", css)
	}
}

func TestSassCEngineCompressed(t *testing.T) {
	css, err := sassc.NewEngine(".a{b:1}", sassc.EngineOptions{Style: sassc.StyleCompressed}).Render()
	if err != nil {
		t.Fatal(err)
	}
	if css != ".a{b:1}\n" {
		t.Errorf("got %q", css)
	}
}

func TestSassCEngineIndented(t *testing.T) {
	css, err := sassc.NewEngine(".a\n  b: 1\n", sassc.EngineOptions{Syntax: "sass"}).Render()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(css, "b: 1") {
		t.Errorf("got %q", css)
	}
}

func TestSassCEngineError(t *testing.T) {
	if _, err := sassc.NewEngine(".a{b:$undef}", sassc.EngineOptions{}).Render(); err == nil {
		t.Error("expected error")
	}
}

// TestMinimaArithmeticAndColor pins the exact constructs from minima's
// _base.scss that a stale go-scss aborted on ("Undefined operation 30px * 2")
// or serialized differently. The expected strings are byte-verified against
// dart-sass (sass-embedded 1.102.0), the jekyll-sass-converter 3.1.0 oracle.
func TestMinimaArithmeticAndColor(t *testing.T) {
	src := `$spacing-unit: 30px;
$content-width: 800px;
$brand-color: #2a7ae2;
.a { max-width: calc(#{$content-width} - (#{$spacing-unit} * 2)); }
.b { color: darken($brand-color, 15%); }
.c { padding: ($spacing-unit * 0.5); }`
	res, err := sass.CompileString(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"max-width: calc(800px - 30px * 2);",
		"color: rgb(8.9993518068%, 33.8251498947%, 66.0986874088%);",
		"padding: 15px;",
	} {
		if !strings.Contains(res.CSS, want) {
			t.Errorf("missing %q in\n%s", want, res.CSS)
		}
	}
}

// TestJekyllPath exercises the exact call jekyll-sass-converter makes.
func TestJekyllPath(t *testing.T) {
	src := `$brand: #2a7ae2;
.highlight { color: $brand; a { &:hover { text-decoration: underline; } } }`
	res, err := sass.CompileString(src, &sass.Options{Style: sass.StyleCompressed})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, ".highlight{color:#2a7ae2}") || !strings.Contains(res.CSS, ".highlight a:hover") {
		t.Errorf("jekyll path: %q", res.CSS)
	}
}
