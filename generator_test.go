package generator

import (
	"testing"

	"bytes"
	"golang.org/x/tools/go/loader"
)

func TestImport(t *testing.T) {
	cfg := loader.Config{}
	cfg.Import("fmt")
	cfg.Import("log")

	prog, err := cfg.Load()
	if err != nil {
		t.Fatal(err)
	}

	g := New(prog)
	path, selector := g.Import("fmt")
	if path != "fmt" {
		t.Errorf("expected: \"fmt\", got: %s", path)
	}

	if selector != "fmt" {
		t.Errorf("expected: \"fmt\", got: %s", selector)
	}

	g = New(prog)
	logPath, err := g.ImportWithAlias("log", "fmt") //importing log as fmt
	if err != nil {
		t.Errorf("expected nil, got: %v", err)
	}
	if logPath != "log" {
		t.Errorf("expected: \"log\", got: %s", logPath)
	}

	fmtPath, fmtSelector := g.Import("fmt")
	if path != "fmt" {
		t.Errorf("expected: \"fmt\", got: %s", fmtPath)
	}

	if fmtSelector != "fmt2" {
		t.Errorf("expected: \"fmt2\", got: %s", fmtSelector)
	}
}

func TestImportWithAlias(t *testing.T) {
	cfg := loader.Config{}
	cfg.Import("fmt")

	prog, err := cfg.Load()
	if err != nil {
		t.Fatal(err)
	}

	g := New(prog)
	path, err := g.ImportWithAlias("fmt", "f")
	if err != nil {
		t.Errorf("expected nil, got: %s", err)
	}
	if path != "fmt" {
		t.Errorf("expected: \"fmt\", got: %s", path)
	}
}

func TestBuildTags(t *testing.T) {
	cfg := loader.Config{}
	cfg.Import("log")

	prog, err := cfg.Load()
	if err != nil {
		t.Fatal(err)
	}

	g := New(prog)

	g.packageName = "buildtags"

	g.AddBuildTags("linux", "darwin,!cgo")
	g.AddBuildTags("386")

	buf := new(bytes.Buffer)
	if _, err := g.WriteTo(buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected :=
`// +build linux darwin,!cgo
// +build 386

package buildtags
/*
This code was automatically generated using github.com/gojuno/generator lib.
			Please DO NOT modify.
*/`

	if expected != buf.String() {
		t.Errorf("expected %s, got %s", expected, buf.String())
	}
}
